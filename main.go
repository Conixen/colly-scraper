package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"colly-scraper/internal/model"

	"github.com/gocolly/colly/v2"
)

func main() {
	var books []model.Book

	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
		colly.Async(true),
	)

	// Set rate limiting - 1 request per second with randomized delay
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	// Set user agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	// Visit book detail pages - extract all book data
	c.OnHTML("article.product_page", func(e *colly.HTMLElement) {
		book := model.Book{}

		// Extract title
		book.Title = e.ChildText("h1")

		// Extract price
		book.Price = e.ChildText("p.price_color")

		// Extract category from breadcrumb
		breadcrumbs := e.ChildTexts("ul.breadcrumb li")
		if len(breadcrumbs) >= 3 {
			book.Category = strings.TrimSpace(breadcrumbs[2])
		}

		// Extract UPC and Pages from product information table
		e.ForEach("table.table tr", func(_ int, row *colly.HTMLElement) {
			header := row.ChildText("th")
			value := row.ChildText("td")

			if header == "UPC" {
				book.UPC = value
			} else if header == "Number of pages" {
				// Parse the number (e.g., "352" from "352 pages")
				fmt.Sscanf(value, "%d", &book.Pages)
			}
		})

		books = append(books, book)
		fmt.Printf("Scraped: %s\n", book.Title)
	})

	// Visit all book listing pages - find links to individual books
	c.OnHTML("article.product_pod h3 a", func(e *colly.HTMLElement) {
		bookURL := e.Attr("href")
		e.Request.Visit(e.Request.AbsoluteURL(bookURL))
	})

	// Handle pagination - visit next pages
	c.OnHTML("li.next a", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		e.Request.Visit(e.Request.AbsoluteURL(nextPage))
	})

	// error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to scrape URL: %s, Error: %v\n", r.Request.URL, err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	fmt.Println("Starting scraper...")
	c.Visit("https://books.toscrape.com/")

	// Wait for all requests to complete
	c.Wait()

	// Export to JSON
	fmt.Printf("\nScraped %d books total\n", len(books))
	fmt.Println("Exporting to books.json...")

	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON:", err)
	}

	err = os.WriteFile("books.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write JSON file:", err)
	}

	fmt.Println("Successfully exported to books.json!")
}
