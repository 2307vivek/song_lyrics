package utils

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func CreateColly(async bool, parallelism int, delay time.Duration) *colly.Collector {
	c := colly.NewCollector()
	c.Async = async

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting to ", r.URL)
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallelism,
		RandomDelay: delay,
	})

	return c
}