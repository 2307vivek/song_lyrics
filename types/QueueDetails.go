package types

type QueueDetails struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Node  string `json:"node"`
	Vhost string `json:"vhost"`
	Status QueueStatus `json:"backing_queue_status"`
}

type QueueStatus struct {
	Length int `json:"len"`
}