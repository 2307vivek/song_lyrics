package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/2307vivek/song-lyrics/types"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/2307vivek/song-lyrics/utils/api"
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

	api.AppStatus.Connections.Es = true
	ES = client
}

func CheckSongExists(url string) bool {
	exists := false

	if ExistsInCache(utils.SONG_BLOOM_FILTER_NAME, url) {
		query := fmt.Sprintf(`{
  		"query": {
    		"match": {
      		"song.url.keyword": "%s"
    		}
  		}
		}`, url)
		exists = existsInDB(query)
	}
	return exists
}

func CheckArtistExists(url string) bool {
	exists := false

	if ExistsInCache(utils.SONG_BLOOM_FILTER_NAME, url) {
		query := fmt.Sprintf(`{
  		"query": {
    		"match": {
      		"song.aritst.url.keyword": "%s"
    		}
  		}
		}`, url)
		exists = existsInDB(query)
	}
	return exists
}

func existsInDB(query string) bool {
	res, err := ES.Search(
		ES.Search.WithIndex(os.Getenv("ES_INDEX")),
		ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		fmt.Println("Cannot search for song", err)
	}

	var searchRes types.EsResponse

	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		log.Fatalf("Failure to to parse response body: %s", err)
	}
	return searchRes.Hits.Total.Value > 0
}
