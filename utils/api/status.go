package api

import "github.com/gin-gonic/gin"

type Status struct {
	Connections ConnectionStatus `json:"connections"`
	FailedUrls  []string         `json:"failed_urls"`
	ScrapedLyrics int `json:"scraped_urls"`
	ScrapedArtists int `json:"scraped_artists"`
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
