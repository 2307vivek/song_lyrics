package queue

import amqp "github.com/rabbitmq/amqp091-go"

var SongQ *SongQueue

func CreateSongQueue(url string, queueName string) {
	channel, queue := initRabbitMq(url, queueName)

	SongQ.Channel = channel
	SongQ.Queue = queue
}

type SongQueue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}