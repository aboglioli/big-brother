package events

import (
	"github.com/aboglioli/big-brother/errors"
)

var (
	ErrConnect         = errors.Internal.New("rabbitmq.connect")
	ErrCreateChannel   = errors.Internal.New("rabbitmq.create_channel")
	ErrDeclareExchange = errors.Internal.New("rabbitmq.declare_exchange")
	ErrDeclareQueue    = errors.Internal.New("rabbitmq.declare_queue")
	ErrBindQueue       = errors.Internal.New("rabbitmq.bind_queue")
	ErrMarshal         = errors.Internal.New("rabbitmq.marshal")
	ErrPublish         = errors.Internal.New("rabbitmq.publish")
	ErrConsume         = errors.Internal.New("rabbitmq.consume")
)
