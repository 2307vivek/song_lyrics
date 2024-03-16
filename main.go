package main

import (
	"os"

	"github.com/2307vivek/song-lyrics/database"
	"github.com/2307vivek/song-lyrics/handler"
	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/joho/godotenv"
)

func main() {
	instance := os.Args[1]
	if instance != "lyrics" && instance != "artist" && instance != "artist_links" {
		os.Exit(1)
	}

	err := godotenv.Load(".env")
	utils.FailOnError(err, "Failed to load .env")

	database.ConnectToElasticSearch(os.Getenv("ELASTIC_SEARCH_URL"), os.Getenv("ELASTIC_SEARCH_USERNAME"), os.Getenv("ELASTIC_SERACH_PASSWORD"))
	database.ConnectToRedis(os.Getenv("REDIS_URL"))

	rabbitmqUrl := os.Getenv("RABBIT_MQ_URL")
	queue.ConnectToRabbitMq(rabbitmqUrl)
	defer queue.Conn.Close()

	if instance == "artist" {
		handler.ScrapeArtists()
	} else if instance == "lyrics" {
		handler.ScrapeLyrics()
	} else if instance == "artist_links" {
		handler.ScrapeArtistLinks()
	}
}
