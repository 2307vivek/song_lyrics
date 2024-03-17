package queue

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/2307vivek/song-lyrics/types"
)

func GetQueues(vhost string) ([]*types.QueueDetails, error) {

	queueUrl := os.Getenv("RABBITMQ_ADDRESS") + "/api/queues/" + vhost

	client := &http.Client{}

	req, err := http.NewRequest("GET", queueUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set basic authentication header
	req.SetBasicAuth(os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err) 
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	queueDetails := []*types.QueueDetails{}
	if resp.Status == "200 OK" {
		err := json.Unmarshal(body, &queueDetails)
		if err != nil {
			fmt.Println("Error unmarshaing:", err)
		}
	}
	return queueDetails, err
}
