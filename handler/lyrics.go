package handler

import (
	"fmt"
	"time"

	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeLyrics() {
	songQ := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQ.Channel.Close()

	c := utils.CreateColly(true, 9, 1000*time.Millisecond)

	c.OnHTML("#lyrics", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()
		link = h.Request.AbsoluteURL(link)

		fmt.Println("scraping", link)

		var ly []string
		h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				ly = append(ly, s.Text())
			}
		})
		// lyrics := strings.Join(ly, " ")
		// fmt.Println(lyrics)
	})

	songs := songQ.Consume(false, 10)

	var forever chan struct{}
	go func() {
		for song := range songs {
			songLink := string(song.Body)
			c.Visit(songLink)
			song.Ack(false)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	c.Wait()

	fmt.Println("Waiting for song links.")
	<-forever
}
