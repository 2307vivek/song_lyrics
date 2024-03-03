package main

import (
	"fmt"
	"os"
	//"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	x := os.Args[0]

	fmt.Println(x)
}

func helper() {
	c := colly.NewCollector()
	c.Async = true

	songs := make([]string, 0)
	lyrics := make([]string, 0)

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("visiting", r.URL)
	// })

	c.OnHTML("#alfabetMusicList a.nameMusic[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = e.Request.AbsoluteURL(link)
		songs = append(songs, link)
	})

	c.OnHTML("#lyrics", func(h *colly.HTMLElement) {

		link := h.Request.URL.String()

		fmt.Println("visiting", link)
		lyrics = append(lyrics, link)

		if link == "https://www.vagalume.com.br/eminem/zelda-rap--remix-.html" {
			h.DOM.Contents().Each(func(i int, s *goquery.Selection) {
				if goquery.NodeName(s) == "#text" {
					fmt.Printf(">>> (%d) >>> %s\n", i, s.Text())
				}
			})
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*",
	// 	Parallelism: 20,
	// 	RandomDelay: 2 * time.Second,
	// })

	c.Visit("https://www.vagalume.com.br/eminem/zelda-rap--remix-.html")

	c.Wait()

	// for _, song := range songs {
	// 	c.Visit(song)
	// }
	// c.Wait()

	fmt.Println("songs length", len(songs))
	fmt.Println("lyrics length", len(lyrics))
}
