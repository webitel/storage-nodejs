package broker

import (
	"context"
	"github.com/webitel/storage/model"
	"sync"
)

type LayeredBrokerAMQPLayer interface {
	Broker
	Subscribe(Handler)
}

type LayeredBroker struct {
	API API
	sync.Mutex
	TmpContext context.Context
	callbacks  map[*Request]chan *Response
	amqp       LayeredBrokerAMQPLayer
}

func NewLayeredBroker(amqp LayeredBrokerAMQPLayer, api API) Broker {
	broker := &LayeredBroker{
		TmpContext: context.TODO(),
		callbacks:  make(map[*Request]chan *Response),
		amqp:       amqp,
	}

	amqp.Subscribe(broker.onRequest)

	return broker
}

func (l *LayeredBroker) onRequest(req *Request) []byte {
	switch req.Api {
	case "test":
		return []byte("Success test")
	default:
		return []byte("Not found")
	}
}

func (l *LayeredBroker) Close() {
	l.amqp.Close()
}

func (l *LayeredBroker) RPC(request *Request) (*Response, *model.AppError) {
	return l.amqp.RPC(request)
}
