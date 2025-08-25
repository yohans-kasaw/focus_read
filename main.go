package main

import (
	"encoding/json"
	"focusRead/cli"
	"focusRead/epub"
	"os"
)

func main() {

	file_path := "./test_file/test_n_3.epub"
	e, err := epub.New(file_path)
	if err != nil {
		panic(err)
	}

	cli.Run(e.Texts)
	saveTofile_for_debug(e.Texts)
}

func saveTofile_for_debug(texts []epub.Text){
	file, _ := os.Create("./cache/testing_file.txt")
	jsonData, _ := json.MarshalIndent(texts, "", "  ")
	file.Write(jsonData)
}
