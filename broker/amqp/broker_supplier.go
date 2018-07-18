package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"os"
	"time"
)

const (
	PING_ATTEMPTS = 100
	PING_TIME_SEC = 2

	EXIT_AMQP_OPEN = 101
)

type AMQP struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	settings *model.BrokerSettings
	stop     chan struct{}
}

func NewBrokerSupplier(settings model.BrokerSettings) *AMQP {
	supplier := &AMQP{
		settings: &settings,
		stop:     make(chan struct{}, 1),
	}

	supplier.initConnection()
	return supplier
}

func (self *AMQP) initConnection() {
	var err error
	var ch *amqp.Channel

	for {
		self.conn = setupConnection(*self.settings.ConnectionString)
		ch, err = self.conn.Channel()
		if err == nil {
			break
		}
		mlog.Error(fmt.Sprintf("Failed to open AMQP connection to err:%v", err.Error()))
		time.Sleep(time.Second)
	}

	self.channel = ch
	msgs, err := ch.Consume(
		"cdr-leg-a", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	go func() {

		select {
		case d := <-msgs:
			mlog.Debug(fmt.Sprintf("Received a message: %s", d.Body))
		case <-self.stop:
			mlog.Debug("Channel received stop signal.")
			return
		}
		fmt.Println("Close cdr-leg-a channel")
		self.initConnection()
	}()
}

func (self *AMQP) Close() {
	self.stop <- struct{}{}
	self.channel.Close()
	self.conn.Close()
	mlog.Debug("Broker stopped.")
}

func setupConnection(dial string) *amqp.Connection {
	var conn *amqp.Connection
	var err error
	for i := 0; i < PING_ATTEMPTS; i++ {
		conn, err = amqp.Dial(dial)
		if err == nil {
			break
		}
		time.Sleep(time.Second * PING_TIME_SEC)
	}

	if err != nil {
		mlog.Critical(fmt.Sprintf("Failed to open AMQP connection to err:%v", err.Error()))
		time.Sleep(time.Second)
		os.Exit(EXIT_AMQP_OPEN)
	}

	return conn
}
