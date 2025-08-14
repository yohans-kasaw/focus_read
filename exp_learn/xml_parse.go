package explearn

import (
	"encoding/xml"
	"fmt"
)

const xml_str = `
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:dc="http://purl.org/dc/elements/1.1/"
xmlns:ac="http://purl.org/ac/elements/1.1/"
	>
	<head>
		<title>My Simple Page</title>
	</head>
	<body>
		<dc:h1>Welcome dc!</dc:h1>
		<ac:h1>Welcome ac!</ac:h1>
	</body>
</html>
`

type Book struct {
	XMLName xml.Name `xml:"html"`
	Head    struct {
		Title string `xml:"title"`
	} `xml:"head"`
	Body struct {
		DCHeadline string `xml:"http://purl.org/dc/elements/1.1/ h1"`
		ACHeadline string `xml:"http://purl.org/ac/elements/1.1/ h1"`
	} `xml:"body"`
}

func Xml_parse() {
	var book Book
	err := xml.Unmarshal([]byte(xml_str), &book)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(book.XMLName)
	fmt.Println(book.Head.Title)
	fmt.Println("dc ->", book.Body.DCHeadline)
	fmt.Println("ac ->", book.Body.ACHeadline)
}
