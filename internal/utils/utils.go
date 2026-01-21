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

// HasTabletChanged checks if any important fields of a tablet have changed
func HasTabletChanged(old, new model.Tablet) bool {
	return IsScrapedChanged(old.Name, new.Name) ||
		IsScrapedChanged(old.Price, new.Price) ||
		IsScrapedChanged(old.Description, new.Description) ||
		!slicesEqual(old.Color, new.Color) ||
		!slicesEqual(old.HardDrive, new.HardDrive)
}

// slicesEqual checks if two string slices are equal
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}