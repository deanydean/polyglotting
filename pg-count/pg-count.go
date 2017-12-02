package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"
)

func main() {
	// Work out where we are first
	var cwd, err = os.Getwd()

	if err != nil {
		fmt.Println("Failed to get working directory", err)
		return
	}

	// Get the glots indices
	var glots = GetNewGlotsList()

	// Find everything in the current  directory
	findInDir(cwd, glots)

	// Sort results from highest number to lowest
	sort.Sort(sort.Reverse(glots))

	// Write the results out
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	var totalFiles = int64(0)
	var totalLines = int64(0)
	for i := range glots {
		var index = glots[i]

		if index.fileCount > 0 {
			fmt.Fprintf(w, "%s\tfiles=%d\tlines=%d\n",
				index.sourceType, index.fileCount, index.lineCount)

			// Add to the totals
			totalFiles += index.fileCount
			totalLines += index.lineCount
		}
	}

	w.Flush()
	fmt.Println("====================")
	fmt.Fprintf(w, "Total\t%d files\t%d lines\n", totalFiles, totalLines)
	w.Flush()
}
