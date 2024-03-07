package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeLyrics() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	songQ := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQ.Channel.Close()

	c := utils.CreateColly(false, 20, 2*time.Second)

	c.OnHTML("#lyrics", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()

		fmt.Println("visiting", link)

		var ly []string
		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				ly = append(ly, s.Text())
			}
		})
		lyrics := strings.Join(ly, " ")
		fmt.Println(lyrics)
	})

	songs := songQ.Consume(ctx, false)

	var forever chan struct{}
	go func() {
		for song := range songs {
			songLink := song.Body
			c.Visit(string(songLink))
			song.Ack(false)
		}
	}()
	//c.Wait()

	fmt.Println("Waiting for song links.")
	<-forever
}
