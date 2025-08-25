package main

import (
	"fmt"
	"focusRead/epub"
)

func main() {

	file_path := "./test_file/test_v2.epub"
	e, err := epub.New(file_path)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(e.Texts))
	// for _, v := range arr {
	// 	fmt.Println(v)
	// }
}
