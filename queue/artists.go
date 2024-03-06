package queue

import (
	"context"
	"fmt"
	"github.com/2307vivek/song-lyrics/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateArtistQueue(queueName string) (*ArtistQueue) {
	channel, queue := createChannel(queueName)

	artistQ := &ArtistQueue{channel, queue}
	return artistQ
}

func (queue *ArtistQueue) Publish(ctx context.Context, msg []byte) {
	err := queue.Channel.PublishWithContext(ctx,
		"",
		queue.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	utils.FailOnError(err, fmt.Sprintf("Failed to publish message for queue %s\n", queue.Queue.Name))
}

func (queue *ArtistQueue) Consume(ctx context.Context, autoAck bool) <- chan amqp.Delivery {
	msg, err := queue.Channel.Consume(
		queue.Queue.Name,
		"",    // consumer
		autoAck,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	utils.FailOnError(err, fmt.Sprintf("Failed to consume messages for queue %s\n", queue.Queue.Name))

	return msg
}

type ArtistQueue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}