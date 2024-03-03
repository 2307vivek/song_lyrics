package main

import (
	"os"
	"github.com/2307vivek/song-lyrics/handler"
)

func main() {
	instance := os.Args[1]

	if instance == "lyrics" {
		handler.ScrapeArtists()
	} else if instance == "artists" {
		handler.ScrapeArtists()
	} else {
		os.Exit(1)
	}
}
