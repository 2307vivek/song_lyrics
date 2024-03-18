package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/2307vivek/song-lyrics/database"
	"github.com/2307vivek/song-lyrics/types"
)

var SongLyricStore = make(chan types.SongLyrics, 50)

func SongStore(batchSize int) {

	counter := 0
	songBuffer := bytes.NewBuffer(nil)

	timer := time.NewTimer(20 * time.Second)

	go func() {
		for range timer.C {
			if songBuffer.Len() != 0 {
				flush(songBuffer)
				counter = 0
			}
		}
	}()

	for song := range SongLyricStore {
		timer.Reset(20 * time.Second)

		songJson, err := json.Marshal(song)
		if err != nil {
			fmt.Println("failed to serialize song", song)
		}
		songJson = append(songJson, "\n"...)

		if counter >=  batchSize {
			// flush to db
			flush(songBuffer)
			counter = 0
		} else {
			meta := []byte(fmt.Sprintf(`{ "index" : { }%s`, "\n"))
			songBuffer.Write(meta)
			songBuffer.Write(songJson)
			counter++
		}
	}
}

func flush(buffer *bytes.Buffer) {
	fmt.Printf("\"flushing to db\": %v\n", "flushing to db")
	res, e := database.ES.Bulk(bytes.NewReader(buffer.Bytes()), database.ES.Bulk.WithIndex(os.Getenv("ES_INDEX")))
	if e != nil {
		fmt.Println("failed to insert items", e)
	}
	res.Body.Close()
	buffer.Reset()
}
