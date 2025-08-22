package main

import (
	"fmt"
	"focusRead/epub"
)

func main() {
	test("./test_file/test_v2.epub")

}

func test(file_name string){
	fmt.Println("testing --> ", file_name)
	fmt.Println("---------")

	epub, err := epub.New(file_name)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("result")
	fmt.Println("---------")
	fmt.Println("version: ",epub.Package.Version)
	fmt.Println("Nav file: ", epub.Package.NavigationFile)
	fmt.Println("Title", epub.Package.Title)
	fmt.Println("textFiles", epub.TextFiles[0])
	fmt.Println("NavPoints epub 2", epub.Toc.NavPoints[1].NavLabel)
}
