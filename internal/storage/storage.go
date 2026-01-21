package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"colly-scraper/internal/model"
)

const storageDir = "data"

func LoadBook(url string) (model.Book, error) {
	filename := filepath.Join(storageDir, hashURL(url)+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return model.Book{}, nil
		}
		return model.Book{}, fmt.Errorf("failed to read book file: %w", err)
	}

	var book model.Book
	if err := json.Unmarshal(data, &book); err != nil {
		return model.Book{}, fmt.Errorf("failed to unmarshal book: %w", err)
	}

	return book, nil
}

func SaveBook(book model.Book) error {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	filename := filepath.Join(storageDir, hashURL(book.Pagelink)+".json")

	data, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal book: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write book file: %w", err)
	}

	return nil
}

func hashURL(url string) string {
	hash := ""
	for _, c := range url {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			hash += string(c)
		} else {
			hash += "_"
		}
	}
	if len(hash) > 100 {
		hash = hash[:100]
	}
	return hash
}

func LoadTablet(name string) (model.Tablet, error) {
	filename := filepath.Join(storageDir, "tablet_"+hashURL(name)+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return model.Tablet{}, nil
		}
		return model.Tablet{}, fmt.Errorf("failed to read tablet file: %w", err)
	}

	var tablet model.Tablet
	if err := json.Unmarshal(data, &tablet); err != nil {
		return model.Tablet{}, fmt.Errorf("failed to unmarshal tablet: %w", err)
	}

	return tablet, nil
}

func SaveTablet(tablet model.Tablet) error {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	filename := filepath.Join(storageDir, "tablet_"+hashURL(tablet.Name)+".json")

	data, err := json.MarshalIndent(tablet, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tablet: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write tablet file: %w", err)
	}

	return nil
}
