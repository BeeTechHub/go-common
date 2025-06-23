package utils

import (
	"math"
)

func CalculatePaginatedSkip(page, size int64) int64 {
	skip := page*size - size
	return skip
}

func CalculatePageCount(recordCount int64, pageSize int64) int64 {
	pageCount := float64(recordCount) / float64(pageSize)
	_pageCount := int64(math.Ceil(pageCount))
	return _pageCount
}
