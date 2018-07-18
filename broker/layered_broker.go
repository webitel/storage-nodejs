package broker

import (
	"context"
)

type LayeredBrokerAMQPLayer interface {
	Broker
}

type LayeredBroker struct {
	TmpContext context.Context
	amqp       LayeredBrokerAMQPLayer
}

func NewLayeredBroker(amqp LayeredBrokerAMQPLayer) Broker {
	broker := &LayeredBroker{
		TmpContext: context.TODO(),
		amqp:       amqp,
	}

	return broker
}

func (l *LayeredBroker) Close() {
	l.amqp.Close()
}
