package main

import (
	"fmt"
	"os"

	"github.com/deanydean/polyglotting/pg"
)

// Weightings
var avgWeight = float64(0.000001)
var devWeight = float64(0.00001)
var maxWeight = float64(0.1)
var minWeight = float64(10)

func calcIndexRating(glots pg.GlotIndices) float64 {
	// Get a mean average over the glots
	// TODO review average (should be median?)
	var glotCount = 0
	var glotCounterSum = float64(0)
	var glotMaxLines = float64(-1)
	var glotMinLines = float64(-1)
	for i := range glots {
		if glots[i].LineCount > 0 {
			glotCount++

			if float64(glots[i].LineCount) > glotMaxLines {
				glotMaxLines = float64(glots[i].LineCount)
			}

			if glots[i].LineCount != 0 &&
				(glotMinLines == -1 ||
					float64(glots[i].LineCount) < glotMinLines) {
				glotMinLines = float64(glots[i].LineCount)
			}

			glotCounterSum += float64(glots[i].LineCount)
		}
	}
	var glotAvg = glotCounterSum / float64(glotCount)

	// Get an average deviation rating
	var deviation = (glotMaxLines - glotMinLines) * devWeight

	// Get the number of lines between min, max and rest
	var maxOffset = float64(0)
	var minOffset = float64(0)
	for i := range glots {
		maxOffset += glotMaxLines - float64(glots[i].LineCount)
		minOffset += float64(glots[i].LineCount) - glotMinLines
	}

	fmt.Println("Raw data points: ")
	fmt.Println("glots", glotCount, "avg", glotAvg, "dev", deviation,
		"maxLines", glotMaxLines, "minLines", glotMinLines,
		"maxOff", maxOffset, "minOff", minOffset)

	var calcMaxDev = deviation / (maxOffset * maxWeight)
	var calcMinDev = deviation * (minOffset * minWeight)

	fmt.Println("Calc values:")
	fmt.Println("calcMaxDev", calcMaxDev, "calcMinDev", calcMinDev)

	// Return rating
	return (glotAvg * avgWeight) * (calcMinDev + calcMaxDev)
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

	fmt.Printf("Polyglot rating is %.2f\n", rating)
}
