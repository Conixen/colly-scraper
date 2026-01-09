package utils

import "strings"


func CleanScraped(s string) string {	// cleans up whitespaces if there are any
	return strings.TrimSpace(s)
}