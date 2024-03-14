package handler

import (
	"fmt"

	"github.com/2307vivek/song-lyrics/types"
)

var SongLyricStore = make(chan types.SongLyrics, 50)

func SongStore(count int) {

	index := 0
	songs := []types.SongLyrics{}

	for song := range SongLyricStore {
		if index == count {
			// flush to db
			fmt.Println("FLushed to db")
			index = 0
		} else {
			songs = append(songs, song)
			index++
		}
	}
}
