package main

import (
	"flag"
	"fmt"
	"focusRead/cli"
	"os"
)

const (
	cacheDir     = "./cache"
	pasteDir     = "./paste"
	progressFile = "history.json"
)

func main() {
	pasteOption := flag.Bool("paste", false, "paste content and read")
	flag.Parse()

	ps, err := NewProgressStore()
	if err != nil {
		fmt.Printf("Error loading reading progress: %v\n", err)
	}

	var bookPath string
	var texts []string

	if *pasteOption {
		bookPath, texts, err = handlePasteMode()
	} else {
		bookPath, texts, err = handleNormalMode(ps)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	progress := ps.GetProgress(bookPath)
	endIndex := cli.Run(texts, progress.Index)

	progress.Index = endIndex
	if err := ps.SaveProgress(); err != nil {
		fmt.Printf("Warning: Error saving progress: %v\n", err)
	}
	if err := debugWriteTexts(texts, bookPath); err != nil {
		fmt.Printf("Warning: Failed to write texts to debug cache. You can ignore this in production. Details: %v\n", err)
	}
}
