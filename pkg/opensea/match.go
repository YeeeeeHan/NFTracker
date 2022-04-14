package opensea

import (
	"github.com/schollz/closestmatch"
)

func FindClosestmatch(slug string, n int, wordsToTest []string) []string {
	// Choose a set of bag sizes, more is more accurate but slower
	bagSizes := []int{2, 3, 4}

	// Create a closestmatch object
	cm := closestmatch.New(wordsToTest, bagSizes)

	return cm.ClosestN(slug, n)
}
