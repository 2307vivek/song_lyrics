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

	c := utils.CreateColly(true, 20, 2 * time.Second)

	c.OnHTML("#alfabetMusicList a.nameMusic[href]", func(h *colly.HTMLElement) {
		songLink := h.Attr("href")
		
		fmt.Println("publishing", songLink)
		queue.SongQ.Publish(ctx, []byte(songLink))
	})

	c.Visit("https://www.vagalume.com.br/eminem/")
	c.Wait()
}