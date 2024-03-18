package api

type Status struct {
	Connections ConnectionStatus `json:"connections"`
	FailedUrls  []string         `json:"failed_urls"`
}

type ConnectionStatus struct {
	Redis    bool `json:"redis"`
	Es       bool `json:"es"`
	RabbitMQ bool `json:"rabbit_mq"`
}

var AppStatus Status = Status{}
