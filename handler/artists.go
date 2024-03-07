package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/gocolly/colly"
)

func ScrapeArtists() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	songQ := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQ.Channel.Close()
	artistQ := queue.CreateArtistQueue(utils.ARTIST_QUEUE_NAME)
	defer artistQ.Channel.Close()

	c := utils.CreateColly(false, 20, 2 * time.Second)

	c.OnHTML("#alfabetMusicList a.nameMusic[href]", func(h *colly.HTMLElement) {
		songLink := h.Attr("href")
		
		fmt.Println("publishing", songLink)
		
		songQ.Publish(ctx, []byte(songLink))
	})

	artists := artistQ.Consume(ctx, false)

	var forever chan struct{}
	go func() {
		for artist := range artists {
			artistLink := artist.Body
			c.Visit(string(artistLink))
			artist.Ack(false)
		}
	}()

	fmt.Println("Waiting for artist links.")
	<- forever
}