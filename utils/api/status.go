package api

import (
	"sync"

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

var AppStatus Status = Status{}

func InitAppStatus() {
	router := gin.Default()

	router.GET("/status", getStatus)
	router.Run(":8080")
}

func getStatus(c *gin.Context) {
	c.IndentedJSON(200, AppStatus)
}

var mutex *sync.RWMutex = &sync.RWMutex{}

func IncrementCountArtist() {
	mutex.Lock()
	AppStatus.ScrapedArtists++
	mutex.Unlock()
}

func IncrementCountLyrics() {
	mutex.Lock()
	AppStatus.ScrapedLyrics++
	mutex.Unlock()
}
