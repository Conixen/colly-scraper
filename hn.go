package main

import (
	"fmt"
	"strings"
	"time"

	"colly-scraper/internal/model"
	"colly-scraper/internal/utils"
	"colly-scraper/internal/storage"

	"github.com/gocolly/colly/v2"
)

func main() {

	scrapeUrl := "https://books.toscrape.com/"

	c := colly.NewCollector(colly.AllowedDomains("books.toscrape.com"))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	var newBook model.Book

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("Accept-Language", "en-US;q=0.9") 		// if webbsite supports multiple languages
		fmt.Printf("Start scraping %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Something went wrong while scraping %s: %v\n", r.Request.URL, err)
	})

	c.OnHTML("article.product_page", func(e *colly.HTMLElement) {
		newBook = model.Book{}

		title := e.ChildText("h1")
		fmt.Printf("Book Title found: %s\n", title)
		newBook.Title = title

		category := e.ChildText("ul.breadcrumb li:nth-child(3)")
		newBook.Category = strings.TrimSpace(category)
		fmt.Printf("Book Category found: %s\n", newBook.Category)

		img := e.ChildAttr("div.item.active img", "src")
		newBook.IMGURL = img
		fmt.Printf("Book Image URL found: %s\n", img)

		e.ForEach("table.table tr", func(_ int, row *colly.HTMLElement) {
			header := row.ChildText("th")
			value := utils.CleanScraped(row.ChildText("td"))

			if header == "UPC" {
				newBook.UPC = value
				fmt.Printf("Book UPC found: %s\n", value)
			}

			if header == "Price (excl. tax)" {
				newBook.Price = value
				fmt.Printf("Book Price found: %s\n", value)
			}

			if header == "Availability" {
				newBook.InStock = value
				fmt.Printf("Book Availability found: %s\n", value)
			}
		})

		newBook.Pagelink = e.Request.URL.String()
		fmt.Printf("Book scraped successfully: %s\n", newBook.Title)
	})

	c.OnScraped(func(r *colly.Response) {
		oldBook, err := storage.LoadBook(r.Request.URL.String())
		if err != nil {
			fmt.Printf("Warning: failed to load old book data: %v\n", err)
		}

		if utils.HasBookChanged(oldBook, newBook) {
			fmt.Printf("Changes detected! Updating book: %s\n", newBook.Title)

			if err := storage.SaveBook(newBook); err != nil {
				fmt.Printf("Error: failed to save book: %v\n", err)
			} else {
				fmt.Printf("Book saved successfully: %s\n", newBook.Title)
			}
		} else {
			fmt.Printf("No changes detected for: %s - skipping update\n", newBook.Title)
		}
	})
	c.Visit(scrapeUrl)

	// for each book link on the page
	c.OnHTML("article.product_pod h3 a", func(e *colly.HTMLElement) {
	 	bookURL := e.Attr("href")
	 	e.Request.Visit(e.Request.AbsoluteURL(bookURL))
  	})

	// for next page
	c.OnHTML("li.next a", func(e *colly.HTMLElement) {
	 	nextPage := e.Attr("href")
	 	e.Request.Visit(e.Request.AbsoluteURL(nextPage))
	})
}