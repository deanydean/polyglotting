package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/deanydean/polyglotting/pg"
)

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

	// Sort results from highest number to lowest
	sort.Sort(sort.Reverse(glots))

	// Write the results out
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	var totalFiles = int64(0)
	var totalLines = int64(0)
	for i := range glots {
		var index = glots[i]

		if index.FileCount > 0 {
			fmt.Fprintf(w, "%s\tfiles=%d\tlines=%d\n",
				index.SourceType, index.FileCount, index.LineCount)

			// Add to the totals
			totalFiles += index.FileCount
			totalLines += index.LineCount
		}
	}

	w.Flush()
	fmt.Println("====================")
	fmt.Fprintf(w, "Total\t%d files\t%d lines\n", totalFiles, totalLines)
	w.Flush()
}
