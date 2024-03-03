package queue

import (
	"fmt"

	"github.com/2307vivek/song-lyrics/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func connectToRabbitMq(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to rabbitmq.")

	return conn
}

func initRabbitMq(url string, queueName string) (*amqp.Channel, *amqp.Queue, *amqp.Connection) {
	conn := connectToRabbitMq(url)

	channel, err := conn.Channel()
	utils.FailOnError(err, fmt.Sprintf("Failed to create channel %s\n", queueName))

	queue, err := channel.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, fmt.Sprintf("Failed to create queue %s\n", queueName))

	return channel, &queue, conn
}