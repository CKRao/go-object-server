package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatalf("amqp.Dial error : %s , url %s \n", err, url)
	}
	channel, err := conn.Channel()

	if err != nil {
		log.Fatalf("conn.Channel error : %s  \n", err)
	}

	queue, err := channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("channel.QueueDeclare error : %s  \n", err)
	}

	mq := new(RabbitMQ)

	mq.channel = channel
	mq.Name = queue.Name

	return mq
}

//绑定交换机
func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("QueueBind error : %s  \n", err)
	}

	q.exchange = exchange
}

//发送消息
func (q *RabbitMQ) Send(queue string, body interface{}) {
	bytes, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Send -> json.Marshal error : %s , queue : %s \n", err, queue)
	}

	err = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{ReplyTo: q.Name, Body: bytes},
	)

	if err != nil {
		log.Fatalf("channel.Publish error : %s , queue : %s \n", err, q.Name)
	}
}

//发布消息
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	bytes, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Publish -> json.Marshal error : %s , exchange : %s \n", err, exchange)
	}

	err = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{ReplyTo: q.Name, Body: bytes},
	)

	if err != nil {
		log.Fatalf("channel.Publish error : %s , exchange : %s \n", err, exchange)
	}
}

//消费消息
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	deliveries, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("q.channel.Consume : %s , queue : %s \n", err, q.Name)
	}

	return deliveries
}

//关闭消息队列
func (q *RabbitMQ) Close() {
	q.channel.Close()
}
