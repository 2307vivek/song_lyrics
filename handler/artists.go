package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/2307vivek/song-lyrics/database"
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

	vhost := os.Getenv("VHOST")

	c := utils.CreateColly(false, 20, 2*time.Second)

	c.OnHTML("body", func(h *colly.HTMLElement) {
		artistLink := h.Request.URL.String()
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
		database.AddToCache(utils.ARTIST_BLOOM_FILTER_NAME, artistLink)
	})

	artists := artistQ.Consume(false, 1)

	var forever chan struct{}
	go func() {
		for artist := range artists {
			artistLink := string(artist.Body)

			for {
				queueLength := getQueueLength(vhost)
				fmt.Printf("songQueueLen: %v\n", queueLength)
				if queueLength < 8500 {
					if !database.Exists(utils.ARTIST_BLOOM_FILTER_NAME, artistLink) {
						fmt.Printf("artistLink: %v\n", artistLink)
						c.Visit(artistLink)
					}
					artist.Ack(false)
					break
				} else {
					timer := time.NewTimer(2 * time.Minute)
					fmt.Println("waiting to consume songs")
					<-timer.C
				}
			}
		}
	}()

	fmt.Println("Waiting for artist links.")
	<-forever
}

func getQueueLength(vhost string) int {
	queues, err := queue.GetQueues(vhost)

	if err != nil {
		fmt.Println("cannot get details", err)
	}
	length := 0
	for _, queue := range queues {
		length += queue.Status.Length
	}
	return length
}
