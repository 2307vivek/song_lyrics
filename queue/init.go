package queue

import (
	"github.com/2307vivek/song-lyrics/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func connectToRabbitMq(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to rabbitmq.")
	defer conn.Close()

	return conn
}

func initRabbitMq(url string, queueName string) (*amqp.Channel, *amqp.Queue) {
	conn := connectToRabbitMq(url)

	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to create channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to create queue")

	return channel, &queue
}