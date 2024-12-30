package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

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

	// 创建队列
	q, err := ch.QueueDeclare("go_q1",
		true,  // durable
		false, //autoDelete
		false, //exclusive
		false, //noWait
		nil,   //args
	)

	if err != nil {
		panic(err)
	}

	// 获取消息
	go consume("c1", conn, q.Name)
	go consume("c2", conn, q.Name)

	// 发送消息
	i := 0
	for {
		i++
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
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

func consume(consumer string, conn *amqp.Connection, q string) {
	// 建立连接
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

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
