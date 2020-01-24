package amqp

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/webitel/storage/broker"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
	"time"
)

const (
	PING_ATTEMPTS = 100
	PING_TIME_SEC = 2

	EXIT_AMQP_OPEN = 101

	QUEUE_RPC_NAME = "storage.rpc"

	CONTEXT_TYPE_JSON = "application/json"
)

type session struct {
	*amqp.Connection
	*amqp.Channel
}

func (s *session) Close() {
	if s.Channel != nil {
		s.Channel.Close()
	}

	if s.Connection != nil {
		s.Connection.Close()
	}
}

type AMQP struct {
	settings *model.BrokerSettings
	ctx      context.Context
	done     context.CancelFunc
	handler  broker.Handler

	stoppedSubscriber chan struct{}
}

func NewBrokerSupplier(settings model.BrokerSettings) *AMQP {
	supplier := &AMQP{
		settings:          &settings,
		stoppedSubscriber: make(chan struct{}, 1),
	}
	supplier.ctx, supplier.done = context.WithCancel(context.Background())

	return supplier
}

func (a *AMQP) redial() chan chan session {
	sessions := make(chan chan session)

	go func() {

		rec := make(chan bool, 1)
		sess := make(chan session)
		var s *session

		defer func() {
			wlog.Debug("Shutting down session factory")
			close(sessions)
			if s != nil {
				s.Close()
			}

			close(sess)

		}()

		for {
			select {
			case sessions <- sess:

			case <-a.ctx.Done():
				return
			case <-rec:

			}

			conn, err := amqp.Dial(*a.settings.ConnectionString)
			if err != nil {
				wlog.Error(fmt.Sprintf("cannot (re)dial: %v: %s", err, *a.settings.ConnectionString))
				time.Sleep(time.Second)
				rec <- true
				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				wlog.Error(fmt.Sprintf("cannot create channel: %v", err))
				time.Sleep(time.Second)
				rec <- true
				continue
			}

			s = &session{conn, ch}

			select {
			case sess <- *s:
			case <-a.ctx.Done():
				return
			}
		}

	}()

	return sessions
}

func (a *AMQP) Subscribe(h broker.Handler) {
	a.handler = h
}

func (a *AMQP) subscribe(sessions chan chan session) {
	var deliveries <-chan amqp.Delivery
	var err error
	var sub session
	var ok bool

	defer func() {
		wlog.Debug("Subscriber: finished.")
		a.stoppedSubscriber <- struct{}{}
	}()

	for sess := range sessions {
		sub, ok = <-sess
		if !ok {
			return
		}

		if _, err = sub.QueueDeclare(QUEUE_RPC_NAME, false, false, false, false, nil); err != nil {
			wlog.Error(fmt.Sprintf("Subscriber: cannot consume from exclusive queue: %q, %v", QUEUE_RPC_NAME, err))
			continue
		}

		deliveries, err = sub.Consume(QUEUE_RPC_NAME, "", false, true, false, false, nil)
		if err != nil {
			wlog.Error(fmt.Sprintf("Subscriber: cannot consume from: %q, %v", QUEUE_RPC_NAME, err))
			continue
		}

		wlog.Debug("Subscriber: listen messages ...")

		for msg := range deliveries {
			wlog.Debug(fmt.Sprintf("Receive %d bytes: %s", len(msg.Body), string(msg.Body)))
			err, req := broker.NewRequestFromBytes(msg.Body)
			if err != nil {
				wlog.Error(fmt.Sprintf("Cannot parse request, error: %s", err.Error()))
				sub.Ack(msg.DeliveryTag, false)
				//TODO bad request
			} else if a.handler != nil {
				go func(msg amqp.Delivery) {
					res := a.handler(&req)
					fmt.Println(string(res))

					sub.Publish(msg.Exchange, msg.RoutingKey, false, false, amqp.Publishing{
						CorrelationId: msg.CorrelationId,
						ReplyTo:       msg.ReplyTo,
						Body:          res,
					})

				}(msg)

				sub.Ack(msg.DeliveryTag, false)
			}

		}
	}
}

func (self *AMQP) Close() {
	wlog.Info("Stopping broker")
	self.done()
	<-self.stoppedSubscriber
	wlog.Info("Broker stopped.")
}

func (self *AMQP) RPC(request *broker.Request) (*broker.Response, *model.AppError) {
	//TODO
	return nil, nil
}
