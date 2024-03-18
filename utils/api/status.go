package api

import "github.com/2307vivek/song-lyrics/types"

type Status struct {
	Connections   ConnectionStatus   `json:"connections"`
	FailedUrls    []string           `json:"failed_urls"`
	FailedInserts []types.SongLyrics `json:"failed_inserts"`
}

type ConnectionStatus struct {
	Redis    bool `json:"redis"`
	Es       bool `json:"es"`
	RabbitMQ bool `json:"rabbit_mq"`
}

var AppStatus Status = Status{}
