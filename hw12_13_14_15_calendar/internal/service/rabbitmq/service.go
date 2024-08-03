package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
)

type Rabbit struct {
	connection *amqp.Connection
	ch         *amqp.Channel
	se         *config.RabbitConf
}

func NewRabbit(
	con *amqp.Connection,
	ch *amqp.Channel,
	queue *config.RabbitConf,
) *Rabbit {
	return &Rabbit{
		connection: con,
		ch:         ch,
		se:         queue,
	}
}

func (r *Rabbit) InitialQueue() error {
	var err error
	if err := r.ch.ExchangeDeclare(
		r.se.Producer.Exchange.Name,
		r.se.Producer.Exchange.Type,
		r.se.Producer.Exchange.Durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %s: %w", r.se.Producer.Exchange.Name, err)
	}

	if err = r.ch.ExchangeBind(
		r.se.Producer.Exchange.Name,
		r.se.Producer.RoutingKey,
		r.se.Producer.Exchange.Name,
		false,
		map[string]interface{}{},
	); err != nil {
		return fmt.Errorf("exchange bind: %s: %w", r.se.Producer.Queue.Name, err)
	}

	if _, err = r.ch.QueueDeclare(
		r.se.Producer.Queue.Name,
		r.se.Producer.Queue.Durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue declare: %s: %w", r.se.Producer.Queue.Name, err)
	}

	if err = r.ch.QueueBind(
		r.se.Producer.Queue.Name,
		r.se.Producer.RoutingKey,
		r.se.Producer.Exchange.Name,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue bind: %s: %w", r.se.Producer.Queue.Name, err)
	}

	return nil
}

func (r *Rabbit) PublishMessages(evtss []storage.Event) error {
	for _, event := range evtss {
		res, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal event: %w", err)
		}

		if err := r.ch.Publish(
			r.se.Producer.Exchange.Name,
			r.se.Producer.Queue.Name,
			false,
			false,
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            res,
				DeliveryMode:    amqp.Transient,
				Priority:        0,
			},
		); err != nil {
			return fmt.Errorf("exchange publish msg: %w", err)
		}
	}

	return nil
}

func (r *Rabbit) Consume() (<-chan amqp.Delivery, error) {
	messageChannel, err := r.ch.Consume(
		r.se.Consumer.Queue,
		r.se.Consumer.Consumer,
		r.se.Consumer.AutoAck,
		r.se.Consumer.Exclusive,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("start consume: %w", err)
	}

	return messageChannel, nil
}
