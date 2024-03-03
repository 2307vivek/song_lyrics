package queue

import amqp "github.com/rabbitmq/amqp091-go"

var ArtistQ *ArtistQueue

func CreateArtistQueue(url string, queueName string) (*amqp.Connection, *amqp.Channel) {
	channel, queue, conn := initRabbitMq(url, queueName)

	ArtistQ = &ArtistQueue{channel, queue}

	return conn, channel
}

type ArtistQueue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}