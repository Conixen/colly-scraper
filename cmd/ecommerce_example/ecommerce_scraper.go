package main

import (
	"fmt"
	"time"

	"colly-scraper/internal/model"

	"github.com/gocolly/colly/v2"
)

func main() {
	scrapeUrl := "https://webscraper.io/test-sites"

	c = colly.NewCollector(colly.AllowedDomains("webscraper.io/test-sites"))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	var tablet model.Tablet

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("Accept-Language", "en-US;q=0.9") 		// if webbsite supports multiple languages
		fmt.Printf("Start scraping %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Something went wrong while scraping %s: %v\n", r.Request.URL, err)
	})

	c.OnHTML()
}
