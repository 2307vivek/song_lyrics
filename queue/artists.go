package queue

import amqp "github.com/rabbitmq/amqp091-go"

var ArtistQ *ArtistQueue

func CreateArtistQueue(url string, queueName string) {
	channel, queue := initRabbitMq(url, queueName)

	ArtistQ.Channel = channel
	ArtistQ.Queue = queue
}

type ArtistQueue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}