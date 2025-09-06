package main

import (
	"encoding/json"
	"fmt"
	"focusRead/cli"
	"focusRead/epub"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <path/to/book.epub>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	bookPath := os.Args[1]
	fileInfo, err := os.Stat(bookPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: The specified file was not found: %s\n", bookPath)
		} else {
			fmt.Printf("Error accessing file: %v\n", err)
		}
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		fmt.Printf("Error: The provided path is a directory. Please provide a path to an EPUB file.\n")
		os.Exit(1)
	}

	e, err := epub.New(bookPath)
	if err != nil {
		fmt.Printf("Error creating EPUB reader: Failed to parse EPUB file. Details: %v\n", err)
		os.Exit(1)
	}

	cli.Run(e.Texts)

	fileName := strings.TrimSuffix(
		filepath.Base(bookPath),
		filepath.Ext(bookPath),
	)

	if err := debugWriteTexts(e.Texts, fileName); err != nil {
		log.Printf("Warning: Failed to write texts to debug cache. You can ignore this in production. Details: %v\n", err)
	}
}

func debugWriteTexts(texts []epub.Text, fileName string) error {
	if err := os.MkdirAll("./cache/", 0755); err != nil {
		return fmt.Errorf("error creating directory: %w\n", err)
	}

	file, err := os.Create("./cache/" + fileName + ".txt")
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
