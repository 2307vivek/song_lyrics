package queue

import (
	"fmt"

	"github.com/2307vivek/song-lyrics/utils"
	"github.com/2307vivek/song-lyrics/utils/api"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection

func ConnectToRabbitMq(url string) {
	c, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to rabbitmq.")

	Conn = c

	api.AppStatus.Connections.RabbitMQ = true
}

func createChannel(queueName string) (*amqp.Channel, *amqp.Queue) {

	channel, err := Conn.Channel()
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

	return channel, &queue
}