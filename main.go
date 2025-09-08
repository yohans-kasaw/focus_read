package main

import (
	"encoding/json"
	"fmt"
	"focusRead/cli"
	"focusRead/epub"
	"os"
	"path/filepath"
	"strings"
)

const (
	cacheDir    = "./cache"
	historyFile = "history.json"
)

func main() {
	pm, err := NewProgressManager()
	if len(os.Args) >= 2 {
		path := os.Args[1]
		fileInfo, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Error: The specified file was not found: %s\n", path)
			} else {
				fmt.Printf("Error accessing file: %v\n", err)
			}
			os.Exit(1)
		}

		if fileInfo.IsDir() {
			fmt.Printf("Error: The provided path is a directory. Please provide a path to an EPUB file.\n")
			os.Exit(1)
		}

		pm.SetCurrent(path)
	}

	if err != nil {
		fmt.Printf("Error loading reading progress %v\n", err)
	}

	if len(pm.Progresses) == 0{
		fmt.Printf("Error: No history of reading, please provide book path like the following.\n")
		fmt.Printf("usage: %s <path/to/book.epub>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	path := pm.Progresses[pm.Index].Path
	fmt.Println(path, pm.Index, path)
	e, err := epub.New(path)
	if err != nil {
		fmt.Printf("Error creating EPUB reader: Failed to parse EPUB file. Details: %v\n", err)
		os.Exit(1)
	}

	end_index := cli.Run(e.Texts, pm.Progresses[pm.Index].Index)
	pm.Progresses[pm.Index].Index = end_index
	if err := pm.SaveProgress(); err != nil {
		fmt.Printf("Warning: Error saving progress %v\n", err)
	}

	if err := debugWriteTexts(e.Texts, path); err != nil {
		fmt.Printf("Warning: Failed to write texts to debug cache. You can ignore this in production. Details: %v\n", err)
	}
}

func debugWriteTexts(texts []epub.Text, path string) error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w\n", err)
	}

	fileName := strings.TrimSuffix(
		filepath.Base(path),
		filepath.Ext(path),
	)

	file, err := os.Create(filepath.Join(cacheDir, fileName) + ".txt")
	if err != nil {
		return fmt.Errorf("error creating file: %w\n", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(texts, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marsialing : %w\n", err)
	}

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("Error writing file: %w\n", err)
	}

	return nil
}
