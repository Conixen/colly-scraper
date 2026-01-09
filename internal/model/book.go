package model

type Book struct {
	ID		   int	// gets a unique ID
	Title	   string
	Category   string	// (only books on website) travel, mystery, historical fiction ext
	IMGURL     string
	UPC		   string	// universal product code
	Price      string
	InStock    string 	// availability
	Pagelink   string
	// Author     string
	// Pages	   int
}
