package database

import (
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func ConnectToElasticSearch() {
	client, err := elasticsearch.NewClient()
	utils.FailOnError(err, "Failed to connect to elastic search")

	ES = client
}
