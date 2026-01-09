package model

type Book struct {
	Title	   string
	Author     string
	Pages	   int
	Category   string	// (only books on website) travel, mystery, historical fiction ext
	UPC		   string	// universal product code
	Price      string
}
