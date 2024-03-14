package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/types"
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

	c := utils.CreateColly(false, 20, 2*time.Second)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		fmt.Println("body")
		artist := types.Artist{}

		h.ForEach("#artHeaderTitle .darkBG a", func(i int, e *colly.HTMLElement) {
			artist.Name = e.Text
			artist.Url = h.Request.AbsoluteURL(e.Attr("href"))
		})

		h.ForEach("#artHeaderImg img", func(i int, e *colly.HTMLElement) {
			artist.PicUrl = h.Request.AbsoluteURL(e.Attr("src"))
		})

		h.ForEach("#alfabetMusicList a.nameMusic[href]", func(i int, e *colly.HTMLElement) {
			songName := e.Text
			songLink := e.Attr("href")
			songLink = h.Request.AbsoluteURL(songLink)

			song := types.Song{
				Name:   songName,
				Url:    songLink,
				Artist: artist,
			}

			j, err := json.Marshal(song)
			if err == nil {
				songQ.Publish(ctx, j)
			}
		})
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
	<-forever
}
