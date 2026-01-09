package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	colly.AllowedDomains("https://books.toscrape.com/")

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	// call back if an error occures
	c.OnHTML("*", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})
	
	// error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("failed to scarpe url:", r.Request.URL, "\nresponse:", r, "\nerror:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		println("Visiting", r.URL.String())
	})

	c.Visit("https://books.toscrape.com/")
}
