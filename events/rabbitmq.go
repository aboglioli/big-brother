package events

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/config"
	"github.com/streadway/amqp"
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

func NewRabbitMQ() (Manager, error) {
	config := config.Get()
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		return nil, ErrConnect.M("failed to connect to RabbitMQ with config %s", config.RabbitURL).C("rabbitUrl", config.RabbitURL).Wrap(err)
	}

	return &rabbitMQ{
		conn: conn,
	}, nil
}

func (r *rabbitMQ) Publish(body interface{}, opts *ManagerOptions) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return ErrCreateChannel.M("failed to create channel").Wrap(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		opts.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return ErrDeclareExchange.M("failed to declare exchange %s", opts.Exchange).C("exchange", opts.Exchange).Wrap(err)
	}

	b, err := json.Marshal(body)
	if err != nil {
		return ErrMarshal.M("failed to marshal %v to json", body).Wrap(err)
	}

	err = ch.Publish(
		opts.Exchange,
		opts.Route,
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

func (r *rabbitMQ) Consume(opts *ManagerOptions) (<-chan Message, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, ErrCreateChannel.M("failed to create channel").Wrap(err)
	}

	err = ch.ExchangeDeclare(
		opts.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrDeclareExchange.M("failed to declare exchange %s", opts.Exchange).C("exchange", opts.Exchange).Wrap(err)
	}

	exclusive := true
	if opts.Queue != "" {
		exclusive = false
	}

	q, err := ch.QueueDeclare(
		opts.Queue,
		false,
		false,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, ErrDeclareQueue.M("failed to declare queue %s", opts.Queue).C("queue", opts.Queue).Wrap(err)
	}

	err = ch.QueueBind(
		q.Name,
		opts.Route,
		opts.Exchange,
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
