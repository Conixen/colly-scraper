package main

import (
	"fmt"
	"time"

	"colly-scraper/internal/model"
	"colly-scraper/internal/storage"
	"colly-scraper/internal/utils"

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
		newTablet := model.Tablet{}

		newTablet.Name = utils.CleanScraped(e.ChildText("a.title"))
		fmt.Printf("Tablet found: %s\n", newTablet.Name)

		newTablet.Price = utils.CleanScraped(e.ChildText("h4.price"))
		fmt.Printf("Tablet price: %s\n", newTablet.Price)

		newTablet.Description = utils.CleanScraped(e.ChildText("p.description"))
		fmt.Printf("Tablet description: %s\n", newTablet.Description)

		// Collect all colors
		e.ForEach("div.swatches > button[data-color]", func(_ int, el *colly.HTMLElement) {
			color := utils.CleanScraped(el.Attr("data-color"))
			newTablet.Color = append(newTablet.Color, color)
		})
		fmt.Printf("Tablet colors: %v\n", newTablet.Color)

		// Collect all storage options
		e.ForEach("div.swatches > button[data-storage]", func(_ int, el *colly.HTMLElement) {
			storageOpt := utils.CleanScraped(el.Attr("data-storage"))
			newTablet.HardDrive = append(newTablet.HardDrive, storageOpt)
		})
		fmt.Printf("Tablet storage options: %v\n", newTablet.HardDrive)

		// Load previous tablet data for change detection
		oldTablet, err := storage.LoadTablet(newTablet.Name)
		if err != nil {
			fmt.Printf("Warning: failed to load old tablet data: %v\n", err)
		}

		if utils.HasTabletChanged(oldTablet, newTablet) {
			fmt.Printf("Changes detected! Updating tablet: %s\n", newTablet.Name)

			if err := storage.SaveTablet(newTablet); err != nil {
				fmt.Printf("Error: failed to save tablet: %v\n", err)
			} else {
				fmt.Printf("Tablet saved successfully: %s\n", newTablet.Name)
			}
		} else {
			fmt.Printf("No changes detected for: %s - skipping update\n", newTablet.Name)
		}

		tablets = append(tablets, newTablet)
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
