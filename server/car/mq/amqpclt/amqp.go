package amqpclt

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

// NewPublisher creates a new publisher.
func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	// channel is a virtual connection inside a connection
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	err = declareExchange(ch, exchange)
	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

// Publish publishes a message.
func (p *Publisher) Publish(c context.Context, car *carpb.CarEntity) error {
	b, err := json.Marshal(car)
	if err != nil {
		return fmt.Errorf("cannot marshal car: %v", err)
	}

	return p.ch.Publish(
		p.exchange,
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body: b,
		}, // message
	)
}

type Subscriber struct {
	conn     *amqp.Connection
	exchange string
	logger   *zap.Logger
}

func NewSubscriber(conn *amqp.Connection, exchange string, logger *zap.Logger) (*Subscriber, error) {
	// 先建立exchange，不确定是先执行publisher还是subscriber，所以两边都要建立exchange，后建立的会忽略已存在的exchange
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	defer ch.Close()
	err = declareExchange(ch, exchange)
	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Subscriber{
		conn:     conn,
		exchange: exchange,
		logger:   logger,
	}, nil
}
func (s *Subscriber) SubscribeRaw(context.Context) (<-chan amqp.Delivery, func(), error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot allocate channel: %v", err)
	}
	// 这里不能使用defer关闭channel，在函数退出时会返回msgs，msgs会在外部使用，所以不能在函数退出时关闭channel
	// defer ch.Close()

	// 将close函数返回给调用者，调用者可以在不需要消息时关闭channel
	closeCh := func() {
		err := ch.Close()
		if err != nil {
			s.logger.Error("cannot close channel", zap.Error(err))
		}
	}
	// 创建队列
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  //autoDelete
		false, //exclusive
		false, //noWait
		nil,   //args
	)
	if err != nil {
		return nil, closeCh, fmt.Errorf("cannot declare queue: %v", err)
	}

	// 先删除队列，后关闭channel
	cleanUp := func() {
		_, err := ch.QueueDelete(
			q.Name, // name
			false,  // ifUnused
			false,  // ifEmpty
			false,  // noWait
		)
		if err != nil {
			s.logger.Error("cannot delete queue", zap.String("name", q.Name), zap.Error(err))
		}
		closeCh()
	}
	// 绑定队列
	err = ch.QueueBind(
		q.Name,     // queue
		"",         // key
		s.exchange, // exchange
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot bind queue: %v", err)
	}

	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot consume: %v", err)
	}

	return msgs, cleanUp, nil
}

// 数据转换，将amqp.Delivery转换为carpb.CarEntity
func (s *Subscriber) Subscribe(c context.Context) (chan *carpb.CarEntity, func(), error) {
	msgCh, cleanUp, err := s.SubscribeRaw(c)
	if err != nil {
		return nil, cleanUp, err
	}

	carCh := make(chan *carpb.CarEntity)
	go func() {
		for msg := range msgCh {
			var car carpb.CarEntity
			err := json.Unmarshal(msg.Body, &car)
			if err != nil {
				s.logger.Error("cannot unmarshal car", zap.Error(err))
				continue
			}
			carCh <- &car
		}
		// 上面的for循环会不断从msgCh中获取消息，直到msgCh关闭退出for循环，这时候关闭carCh
		close(carCh)
	}()
	return carCh, cleanUp, nil
}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange, // name
		"fanout", // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}
