package main

import (
	"fmt"
	"os"

	"github.com/deanydean/polyglotting/pg"
)

func calcIndexRating(glots pg.GlotIndices) float64 {
	// Get a polyglot rating for this codebase

	// Get the required info before we start
	var totalLines = int64(0)
	for i := range glots {
		totalLines += glots[i].LineCount
	}

	// Now work out the percentage of code for each glot
	var glotPercents = make([]float64, glots.Len())
	var maxPercent = float64(0)
	for i := range glots {
		var lines = glots[i].LineCount
		var percent = float64(lines) / (float64(totalLines) / 100)

		if percent > maxPercent {
			maxPercent = percent
		}

		glotPercents[i] = percent
	}

	// Work out the point scoring for each glot
	var points = float64(0)
	for i := range glotPercents {
		points += glotPercents[i] / maxPercent
	}

	return points
}

func main() {
	// Work out where we are first
	var cwd, err = os.Getwd()

	if err != nil {
		fmt.Println("Failed to get working directory", err)
		return
	}

	// Get the glots indices
	var glots = pg.GetNewGlotsList()

	// Find everything in the current directory
	pg.FindInDir(cwd, glots)

	var rating = calcIndexRating(glots)

	fmt.Printf("%.3f\n", rating)
}
