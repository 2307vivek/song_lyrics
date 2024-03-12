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

	c := utils.CreateColly(true, 15, 100*time.Millisecond)

	count := 0

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
		count--
		fmt.Println(count)
	})

	songs := songQ.Consume(false, 10)

	var forever chan struct{}
	go func() {
		for song := range songs {
			songLink := string(song.Body)
			c.Visit(songLink)
			time.Sleep(100 * time.Millisecond)
			song.Ack(false)
			count++
		}
	}()
	c.Wait()

	fmt.Println("Waiting for song links.")
	<-forever
}
