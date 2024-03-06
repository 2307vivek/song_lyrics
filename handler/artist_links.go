package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/2307vivek/song-lyrics/queue"
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/gocolly/colly"
)

func ScrapeArtistLinks() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	artistPublishQ := queue.CreateArtistQueue(utils.ARTIST_QUEUE_NAME)
	defer artistPublishQ.Channel.Close()

	c := utils.CreateColly(true, 20, 2 * time.Second)

	c.OnHTML("div.bodyCenter ul.gridList li div a[href]", func(h *colly.HTMLElement) {
		artistLink := h.Attr("href")
		artistLink = h.Request.AbsoluteURL(artistLink)
		
		fmt.Println("publishing", artistLink)
		
		artistPublishQ.Publish(ctx, []byte(artistLink))
	})

	c.Visit("https://www.vagalume.com.br/browse/0-9.html")

	for char := 97; char <= 122; char++ {
		c.Visit("https://www.vagalume.com.br/browse/" + string(rune(char)) + ".html")
	}

	c.Wait()
}

