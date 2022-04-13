package opensea

import (
	"github.com/schollz/closestmatch"
)

func FindClosestmatch(slug string, n int) []string {
	// Take a slice of keys, say band names that are similar
	wordsToTest := []string{
		"doodles-official",
		"azuki",
		"boredapeyachtclub",
		"mutant-ape-yacht-club",
		"official-moar-by-joan-cornella",
		"arcade-land",
		"clonex",
		"los-muertos-world",
		"akumaorigins",
		"beanzofficial",
		"cryptopunks",
		"everai-heroes-duo",
		"official-kreepy-club",
		"frankfrank",
		"froyokittenscollection",
		"vaynersports-pass-vsp",
		"hikarinftofficial",
		"tiny-dinos-eth",
		"rtfkt-mnlth",
		"hakinft-io",
		"sandbox",
	}

	// Choose a set of bag sizes, more is more accurate but slower
	bagSizes := []int{2}

	// Create a closestmatch object
	cm := closestmatch.New(wordsToTest, bagSizes)

	return cm.ClosestN(slug, n)
}
