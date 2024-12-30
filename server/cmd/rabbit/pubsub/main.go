package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

const exchange = "go_ex"

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// channel is a virtual connection inside a connection
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		panic(err)
	}

	// 创建两个队列订阅exchange
	go subscribe(conn, exchange)
	go subscribe(conn, exchange)

	// 发送消息到exchange
	i := 0
	for {
		i++
		err := ch.Publish(
			exchange,
			"",    // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				Body: []byte(fmt.Sprintf("message %d", i)),
			}, // message
		)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}

}

func subscribe(conn *amqp.Connection, ex string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

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
		panic(err)
	}
	defer ch.QueueDelete(
		q.Name, // name
		false,  // ifUnused
		false,  // ifEmpty
		false,  // noWait
	)
	// 绑定队列
	err = ch.QueueBind(
		q.Name, // queue
		"",     // key
		ex,     // exchange
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	// 消费消息
	consume("c", ch, q.Name)
}

func consume(consumer string, ch *amqp.Channel, q string) {

	msgs, err := ch.Consume(
		q,        // queue
		consumer, // consumer
		true,     // autoAck
		false,    // exclusive
		false,    // noLocal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Printf("%s received message: %s\n", consumer, msg.Body)
	}
}
