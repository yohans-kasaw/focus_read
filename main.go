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
	bookPath := "./test_file/test_n_4.epub"
	e, err := epub.New(bookPath)

	if err != nil {
		log.Fatalf("Faled to creat epub reader %v", err)
	}

	cli.Run(e.Texts)

	fileName := strings.TrimSuffix(
		filepath.Base(bookPath),
		filepath.Ext(bookPath),
	)

	if err := debugWriteTexts(e.Texts, fileName); err != nil {
		log.Fatalf("filed to write Texts to debug cache %v\n", err)
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
