package epub

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

const cache_dir = "./cache"

func New(file_path string) (*Epub, error) {
	epub := Epub{
		file_path: file_path,
	}

	_, file_name := path.Split(file_path)
	cache_path := path.Join(cache_dir, file_name) + ".bin"

	if ok, _ := os.Stat(cache_path); ok != nil {
		file, err := os.Open(cache_path)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		decoder := gob.NewDecoder(file)
		var restoredTexts []Text
		err = decoder.Decode(&restoredTexts)
		if err != nil {
			panic(err)
		}
		fmt.Println("resored from cache")
		epub.Texts = restoredTexts
	} else {
		fmt.Println("just parsing instead")
		epub.parse()
		epub.cacheTexts(cache_path)
	}

	return &epub, nil
}

func (e *Epub) parse() {
	r, err := zip.OpenReader(e.file_path)
	if err != nil {
		panic(err)
	}

	e.FileMap = make(map[string]*zip.File)
	e.ZipReader = r

	for _, file := range r.File {
		e.FileMap[file.Name] = file
	}

	if err := e.parseContainer(); err != nil {
		panic(err)
	}

	if err := e.parsePackage(); err != nil {
		panic(err)
	}

	if err := e.parseToc(); err != nil {
		panic(err)
	}

	e.parseFlatNavPoint(e.Toc.NavPoints)

	e.constructTextFiles()

	e.parseEpubToTexts()
}

func (e *Epub) cacheTexts(cache_path string) {
	if len(e.Texts) == 0 {
		fmt.Println("Texts is empty nothing to cache")
		return
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(e.Texts)
	if err != nil {
		fmt.Println("error caching", err)
		return
	}

	file, err := os.Create(cache_path)
	if err != nil {
		fmt.Println("error creating a file", err)
		return
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Successfully wrote gob data to product.gob")
}

func (e *Epub) parseEpubToTexts() {
	arr := make([]Text, 0)

	for _, points := range e.Toc.FlatNavPoints {
		if file, ok := e.FileMap[points.Content.Src]; ok {
			r, err := file.Open()
			if err != nil {
				panic(err)
			}

			doc, err := html.Parse(r)
			if err != nil {
				panic(err)
			}
			extract_and_add(doc, &arr)
		}
	}

	e.Texts = arr
}

func (e *Epub) parseContainer() error {

	file_path, ok := e.FileMap["META-INF/container.xml"]
	if !ok {
		return fmt.Errorf("META-INF/container.xml not found in EPUB")
	}

	return readXml(file_path, &e.Container)
}

func (e *Epub) parsePackage() error {
	if e.Container.Rootfile.Path == "" {
		return fmt.Errorf("rootfile path not found in container.xml")
	}

	package_file, ok := e.FileMap[e.Container.Rootfile.Path]
	if !ok {
		return fmt.Errorf(
			"package file %s not found in EPUB",
			e.Container.Rootfile.Path,
		)
	}

	if err := readXml(package_file, &e.Package); err != nil {
		return err
	}

	for _, item := range e.Package.Manifest.Items {
		if item.Media_type == "application/x-dtbncx+xml" || item.Properties == "nav" {
			nav_path := path.Join(path.Dir(e.Container.Rootfile.Path), item.Href)
			e.Package.NavigationFile = nav_path
			return nil
		}
	}

	return nil
}

func (e *Epub) constructTextFiles() {

	e.TextFiles = make([]string, 0)
	for _, spine := range e.Package.Spine.Items {
		for _, manifest := range e.Package.Manifest.Items {
			if spine.Idref == manifest.Id {
				full_path := path.Join(path.Dir(e.Container.Rootfile.Path), manifest.Href)
				e.TextFiles = append(e.TextFiles, full_path)
			}
		}
	}
}

func (e *Epub) parseToc() error {
	if e.Package.NavigationFile == "" {
		return fmt.Errorf("NavigationFile is null ")
	}

	nav_zip_file, ok := e.FileMap[e.Package.NavigationFile]
	if !ok {
		return fmt.Errorf(
			"nav file %s not found in EPUB",
			e.Container.Rootfile.Path,
		)
	}

	return readXml(nav_zip_file, &e.Toc)
}

func (e *Epub) parseFlatNavPoint(navPoints []NavPoint) {
	if e.Toc.FlatNavPoints == nil {
		e.Toc.FlatNavPoints = make([]NavPoint, 0)
	}

	for _, nv := range navPoints {

		found := false
		for _, existingNv := range e.Toc.FlatNavPoints {
			if e.getFullPath(existingNv.Content.Src) == e.getFullPath(nv.Content.Src) {
				found = true
				break
			}
		}

		if !found {
			nv.Content.Src = e.getFullPath(nv.Content.Src)
			e.Toc.FlatNavPoints = append(e.Toc.FlatNavPoints, nv)
		}

		if len(nv.NavPoints) > 0 {
			e.parseFlatNavPoint(nv.NavPoints)
		}
	}
}

func extract_and_add(node *html.Node, arr *[]Text) {
	if node.Type == html.TextNode && node.Parent != nil {
		text := strings.TrimSpace(node.Data)
		if text != "" && node.Parent.Data != "title" && node.Parent.Data != "head" {
			var textType string
			switch node.Parent.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6", "title":
				textType = "h1"
			case "p", "div", "span", "a", "li", "td":
				textType = "p"
			default:
				textType = "anawn"
			}
			*arr = append(*arr, Text{Text: text, Type: textType})
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extract_and_add(c, arr)
	}
}

func (e *Epub) getFullPath(s string) string {
	file_path, _, _ := strings.Cut(s, "#")
	full_path := path.Join(path.Dir(e.Container.Rootfile.Path), file_path)
	return  full_path
}
