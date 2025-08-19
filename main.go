package main

import (
	"fmt"
	"path"
	"focusRead/epub"
)

func main() {
	abs_file_path := "./test_file/test.epub"

	err, epub := epub.NewEpub(abs_file_path)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(epub.Rootfile.Path)
	root_dir := path.Dir(epub.Rootfile.Path)
	fmt.Println(root_dir)
}
