package events

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/streadway/amqp"
)

// Errors
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

// Message
type rabbitMQMessage struct {
	msg amqp.Delivery
}

func (m *rabbitMQMessage) Body() []byte {
	return m.msg.Body
}

func (m *rabbitMQMessage) Event() Event {
	var e Event
	if err := json.Unmarshal(m.msg.Body, &e); err != nil {
		return Event{}
	}
	return e
}

func (m *rabbitMQMessage) Ack() {
	m.msg.Ack(false)
}

// Manager
type rabbitMQ struct {
	conn *amqp.Connection
}

func NewRabbitMQ() (Bus, error) {
	config := config.Get()
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		return nil, ErrConnect.M("failed to connect to RabbitMQ with config %s", config.RabbitURL).C("rabbitUrl", config.RabbitURL).Wrap(err)
	}

	return &rabbitMQ{
		conn: conn,
	}, nil
}

func (r *rabbitMQ) Publish(body interface{}, topic, route string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return ErrCreateChannel.M("failed to create channel").Wrap(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return ErrDeclareExchange.M("failed to declare exchange %s", topic).C("exchange", topic).Wrap(err)
	}

	b, err := json.Marshal(body)
	if err != nil {
		return ErrMarshal.M("failed to marshal %v to json", body).Wrap(err)
	}

	err = ch.Publish(
		topic,
		route,
		false,
		false,
		amqp.Publishing{
			Body: b,
		},
	)
	if err != nil {
		return ErrPublish.M("failed to publish message %s", string(b)).C("message", string(b)).Wrap(err)
	}

	return nil
}

func (r *rabbitMQ) Consume(topic, route, queue string) (<-chan Message, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, ErrCreateChannel.M("failed to create channel").Wrap(err)
	}

	err = ch.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrDeclareExchange.M("failed to declare exchange %s", topic).C("exchange", topic).Wrap(err)
	}

	exclusive := true
	if queue != "" {
		exclusive = false
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrDeclareQueue.M("failed to declare queue %s", queue).C("queue", queue).Wrap(err)
	}

	err = ch.QueueBind(
		q.Name,
		route,
		topic,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrBindQueue.M("failed to bind queue %s", q.Name).C("queue", q.Name).Wrap(err)
	}

	delivery, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrConsume.M("failed to consume from queue %s", q.Name).C("queue", q.Name).Wrap(err)
	}

	msg := make(chan Message)
	go func() {
		for d := range delivery {
			msg <- &rabbitMQMessage{d}
		}
		close(msg)
	}()

	return msg, nil
}
