package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

//新建匿名队列
func New(url string) *RabbitMQ {
	return NewByName(url, "")
}

//通过名称新建队列
func NewByName(url string, name string) *RabbitMQ {

	err, channel := getChannel(url)

	failOnError(err, "Failed to open a channel")

	queue, err := channel.QueueDeclare(
		name,
		false,
		true,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to bind a queue")

	mq := new(RabbitMQ)

	mq.channel = channel
	mq.Name = queue.Name

	return mq
}

//新建交换机
func NewExchange(url string, exchangeName string) *RabbitMQ {
	err, channel := getChannel(url)

	failOnError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	mq := new(RabbitMQ)

	mq.channel = channel
	mq.exchange = exchangeName
	return mq
}

//获取通道
func getChannel(url string) (error, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return err, channel
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

	//queue = "amq.gen-1VYeg2kS-yS3Um-M0l91cw"
	fmt.Println("replay_to ", queue)
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
		log.Fatalf("channel.Publish error : %s , queue : %s \n", err, queue)
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
