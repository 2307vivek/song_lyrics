package handler

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/2307vivek/song-lyrics/database"
	"github.com/2307vivek/song-lyrics/types"
)

var SongLyricStore = make(chan types.SongLyrics, 50)

func SongStore(count int) {

	counter := 0
	songBuffer := bytes.NewBuffer(nil)

	for song := range SongLyricStore {
		songJson, err := json.Marshal(song)
		if err != nil {
			fmt.Println("failed to serialize song", song)
		}
		songJson = append(songJson, "\n"...)

		if counter == count {
			// flush to db
			fmt.Printf("\"flushing to db\": %v\n", "flushing to db")
			res, e := database.ES.Bulk(bytes.NewReader(songBuffer.Bytes()), database.ES.Bulk.WithIndex("sooong"))
			if e != nil {
				fmt.Println("failed to insert items", e)
			}
			res.Body.Close()
			songBuffer.Reset()
			counter = 0
		} else {
			meta := []byte(fmt.Sprintf(`{ "index" : { }%s`, "\n"))
			songBuffer.Write(meta)
			songBuffer.Write(songJson)
			counter++
		}
	}
}
