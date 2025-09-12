package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"focusRead/epub"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
)

func handlePasteMode() (string, []string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", nil, fmt.Errorf("Error reading from clip beard %v\n", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", nil, fmt.Errorf("There is noting to paste..")
	}
	if len(text) < 100 {
		return "", nil, fmt.Errorf("Pasted text is too short..")
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	fileName := re.ReplaceAllString(text[:40], "-") + "-" + uuid.NewString()[:8] + ".txt"
	path := filepath.Join(pasteDir, fileName)
	if err := WriteToFile(text, path); err != nil {
		return "", nil, fmt.Errorf("processing paste: %w", err)
	}

	texts := strings.Split(text, "\n")
	return path, texts, nil
}

func handleNormalMode(pm *ProgressStore) (string, []string, error) {
	var path string
	if len(flag.Args()) > 0 {
		path = flag.Args()[0]

		err := ValidatePathFile(path)
		if err != nil {
			return "", nil, err
		}

	} else {
		if len(pm.Progresses) == 0 {
			fmt.Println("Error: No reading history found. Please provide a book path.")
			fmt.Printf("usage: %s <path/to/book.epub>\n", filepath.Base(os.Args[0]))
			return "", nil, fmt.Errorf("no history found and no file specified")
		}
		path = pm.Progresses[pm.Index].Path
	}

	texts := make([]string, 0)
	if filepath.Ext(path) == ".epub" {
		e, err := epub.New(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse EPUB file: %w", err)
		}

		fmt.Println(e.Texts[10])
		fmt.Println(e.Texts[10])
		for i := range e.Texts {
			texts = append(texts, e.Texts[i].Text)
		}
	} else {
		text, err := ReadFromFile(path)
		if err != nil {
			return "", nil, err
		}
		texts = strings.Split(text, "\n")
	}

	return path, texts, nil
}

func WriteToFile(text string, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Error creating pasteDir %v\n", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating paste file %v\n", err)
	}

	_, err = file.Write([]byte(text))
	if err != nil {
		return fmt.Errorf("Error writing text file %v\n", err)
	}

	return nil
}

func ValidatePathFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("accessing file: %v", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("the provided path is a directory. Please provide a path to an EPUB file")
	}

	return nil
}

func ReadFromFile(path string) (string, error) {
	if err := ValidatePathFile(path); err != nil {
		return "", err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func debugWriteTexts(texts []string, path string) error {
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
