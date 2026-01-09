package main

import (
	"fmt"
	"strings"
	"time"

	"colly-scraper/internal/model"
	"colly-scraper/internal/utils"

	"github.com/gocolly/colly/v2"
)

func main() {

	scrapeUrl := "https://books.toscrape.com/index.html"

	c := colly.NewCollector(colly.AllowedDomains("https://books.toscrape.com/index.html"))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})
	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("Accept-Language", "en-US;q=0.9") 		// if webbsite supports multiple languages
		fmt.Println("Start scarping %s", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong while scraping %s: %v", r.Request.URL, err)
	})

	c.OnHTML("article.product_page", func(e *colly.HTMLElement) {
		book := model.Book{}

		title := e.ChildText("h1")
		fmt.Printf("Book Title found: %s\n", title)
		book.Title = title

		category := e.ChildText("ul.breadcrumb li:nth-child(3)")
		book.Category = strings.TrimSpace(category)
		fmt.Printf("Book Category found: %s\n", book.Category)

		img := e.ChildAttr("div.item.active img", "src")
		book.IMGURL = img
		fmt.Printf("Book Image URL found: %s\n", img)

		e.ForEach("table.table tr", func(_ int, row *colly.HTMLElement) {
			header := row.ChildText("th")
			value := utils.CleanScraped(row.ChildText("td"))

			if header == "UPC" {
				book.UPC = value
				fmt.Printf("Book UPC found: %s\n", value)
			}

			if header == "Price (excl. tax)" {
				book.Price = value
				fmt.Printf("Book Price found: %s\n", value)
			}

			if header == "Availability" {
				book.InStock = value
				fmt.Printf("Book Availability found: %s\n", value)
			}
		})

		book.Pagelink = e.Request.URL.String()
		fmt.Printf("Book scraped successfully: %s\n", book.Title)
	})
	c.Visit(scrapeUrl)

	// Visit all book listing pages - find links to individual books
	// c.OnHTML("article.product_pod h3 a", func(e *colly.HTMLElement) {
	// 	bookURL := e.Attr("href")
	// 	e.Request.Visit(e.Request.AbsoluteURL(bookURL))
	// })

	// // Handle pagination - visit next pages
	// c.OnHTML("li.next a", func(e *colly.HTMLElement) {
	// 	nextPage := e.Attr("href")
	// 	e.Request.Visit(e.Request.AbsoluteURL(nextPage))
	// })

}
