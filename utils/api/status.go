package api

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Connections    ConnectionStatus `json:"connections"`
	FailedUrls     []string         `json:"failed_urls"`
	ScrapedLyrics  int32     `json:"scraped_urls"`
	ScrapedArtists int32     `json:"scraped_artists"`
}

type ConnectionStatus struct {
	Redis    bool `json:"redis"`
	Es       bool `json:"es"`
	RabbitMQ bool `json:"rabbit_mq"`
}

var countLyrics atomic.Int32
var countArtist atomic.Int32

var AppStatus Status = Status{}

func InitAppStatus() {
	router := gin.Default()

	router.GET("/status", getStatus)
	router.Run(":8080")
}

func getStatus(c *gin.Context) {
	c.IndentedJSON(200, AppStatus)
}

func IncrementCountArtist() {
	countArtist.Add(1)
	AppStatus.ScrapedArtists = countArtist.Load()
}

func IncrementCountLyrics() {
	countLyrics.Add(1)
	AppStatus.ScrapedLyrics = countArtist.Load()
}
