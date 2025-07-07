package utils

import "regexp"

func ContainsHTML(input string) bool {
	// Define a regex to match HTML tags
	re := regexp.MustCompile(`<[^>]+>`)
	return re.MatchString(input)
}
