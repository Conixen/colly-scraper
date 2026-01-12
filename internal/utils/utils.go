package utils

import (
	"strings"

	"colly-scraper/internal/model"
)

// Cleans up whitespaces if there are any
func CleanScraped(s string) string { 
	return strings.TrimSpace(s)
}

// checks if scraped field is different
func IsScrapedChanged(old, new string) bool {
	return strings.TrimSpace(old) != strings.TrimSpace(new)
}

// HasBookChanged checks if any important fields of a book have changed
func HasBookChanged(old, new model.Book) bool { 
	return IsScrapedChanged(old.Title, new.Title) ||
		IsScrapedChanged(old.Price, new.Price) ||
		IsScrapedChanged(old.InStock, new.InStock) ||
		IsScrapedChanged(old.Category, new.Category) ||
		IsScrapedChanged(old.IMGURL, new.IMGURL)
}