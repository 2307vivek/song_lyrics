package database

import (
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/elastic/go-elasticsearch/v7"
)

var ES *elasticsearch.Client

func ConnectToElasticSearch(url string, username string, password string) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
	})
	utils.FailOnError(err, "Failed to connect to elastic search")

	ES = client
}
