package main

import (
	"fmt"
	"time"

	"colly-scraper/internal/model"

	"github.com/gocolly/colly/v2"
)

func main() {
	scrapeUrl := "https://webscraper.io/test-sites/e-commerce/static/computers/tablets"

	c := colly.NewCollector(colly.AllowedDomains("webscraper.io"))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	var tablets []model.Tablet

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Start scraping %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Something went wrong while scraping %s: %v\n", r.Request.URL, err)
	})

	c.OnHTML(".thumbnail", func(e *colly.HTMLElement) {
		tablet := model.Tablet{}

		tablet.Name = e.ChildText("a.title")
		fmt.Printf("Tablet found: %s\n", tablet.Name)

		tablet.Price = e.ChildText("h4.price")
		fmt.Printf("Tablet price: %s\n", tablet.Price)

		tablet.Description = e.ChildText("p.description")
		fmt.Printf("Tablet description: %s\n", tablet.Description)

		// Collect all colors
		e.ForEach("div.swatches > button[data-color]", func(_ int, el *colly.HTMLElement) {
			color := el.Attr("data-color")
			tablet.Color = append(tablet.Color, color)
		})
		fmt.Printf("Tablet colors: %v\n", tablet.Color)

		// Collect all storage options
		e.ForEach("div.swatches > button[data-storage]", func(_ int, el *colly.HTMLElement) {
			storage := el.Attr("data-storage")
			tablet.HardDrive = append(tablet.HardDrive, storage)
		})
		fmt.Printf("Tablet storage options: %v\n", tablet.HardDrive)

		tablets = append(tablets, tablet)
		fmt.Println("---")
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("\nFinished scraping. Found %d tablets.\n", len(tablets))
		for i, t := range tablets {
			fmt.Printf("%d. %s - %s\n", i+1, t.Name, t.Price)
		}
	})

	c.Visit(scrapeUrl)
}
