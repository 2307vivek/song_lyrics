package queue

import (
	"context"

	"github.com/2307vivek/song-lyrics/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var SongQ *SongQueue

func CreateSongQueue(url string, queueName string) (*amqp.Connection, *amqp.Channel) {
	channel, queue, conn := initRabbitMq(url, queueName)

	SongQ = &SongQueue{channel, queue}
	return conn, channel
}

func (queue *SongQueue) Publish(ctx context.Context, msg []byte) {
	err := queue.Channel.PublishWithContext(ctx, 
		"",
		queue.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,
		},
	)
	utils.FailOnError(err, "Failed to publish message")
}

type SongQueue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}