package handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	//"github.com/2307vivek/song-lyrics/database"
	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/types"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeLyrics() {
	songQ := queue.CreateSongQueue(utils.SONG_QUEUE_NAME)
	defer songQ.Channel.Close()

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

		songLyrics := SongLyrics{
			lyric: lyrics,
			song:  songMap[link],
		}

		fmt.Println(songLyrics.lyric)
		fmt.Println(songLyrics.song.Name)

		//fmt.Println(lyrics)
		//database.AddToCache(utils.SONG_BLOOM_FILTER_NAME, songLyrics.song.Artist.Name + ":" + songLyrics.song.Name)
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
				c.Visit(songLink)
				time.Sleep(100 * time.Millisecond)
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

type SongLyrics struct {
	lyric string
	song  types.Song
}
