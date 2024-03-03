package main

import (
	"os"

	"github.com/2307vivek/song-lyrics/handler"
	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/joho/godotenv"
)

func main() {
	instance := os.Args[1]
	if instance != "lyrics" && instance != "artist" {
		os.Exit(1)
	}

	err := godotenv.Load(".env")
	utils.FailOnError(err, "Failed to load .env")

	rabbitmqUrl := os.Getenv("RABBIT_MQ_URL")
	songQueueName := os.Getenv("SONG_QUEUE_NAME")
	artistQueueName := os.Getenv("ARTIST_QUEUE_NAME")

	songConnection, songChannel := queue.CreateSongQueue(rabbitmqUrl, songQueueName)
	defer songConnection.Close()
	defer songChannel.Close()

	artistConnection, artistChannel := queue.CreateArtistQueue(rabbitmqUrl, artistQueueName)
	defer artistConnection.Close()
	defer artistChannel.Close()

	if instance == "artist" {
		handler.ScrapeArtists()
	} else if instance == "lyrics" {
		handler.ScrapeLyrics()
	}
}
