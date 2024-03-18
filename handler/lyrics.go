package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/2307vivek/song-lyrics/database"
	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/types"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeLyrics() {

	go SongStore(100)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	songQ := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQ.Channel.Close()

	songQPublish := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQPublish.Channel.Close()

	songMap := make(map[string]types.Song)

	c := utils.CreateColly(true, 15, 100*time.Millisecond)

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
		lyrics := strings.Join(ly, " ")

		songLyrics := types.SongLyrics{
			Lyric: lyrics,
			Song:  songMap[link],
		}

		delete(songMap, link)

		SongLyricStore <- songLyrics

		database.AddToCache(utils.SONG_BLOOM_FILTER_NAME, songLyrics.Song.Name+songLyrics.Song.Artist.Name)
	})

	c.OnError(func(r *colly.Response, err error) {
		link := r.Request.URL.String()

		song := songMap[link]

		s, e := json.Marshal(song)

		if e != nil {
			songQPublish.Publish(ctx, s)	
		}
	})

	songs := songQ.Consume(false, 10)

	var forever chan struct{}
	go func() {
		for song := range songs {
			var s types.Song
			err := json.Unmarshal(song.Body, &s)

			if err == nil {
				songLink := s.Url
				songMap[songLink] = s
				if !database.Exists(utils.SONG_BLOOM_FILTER_NAME, s.Name+s.Artist.Name) {
					c.Visit(songLink)
					time.Sleep(100 * time.Millisecond)
				}
				song.Ack(false)
			} else {
				fmt.Println(err)
			}
		}
	}()
	c.Wait()

	fmt.Println("Waiting for song links.")
	<-forever
}
