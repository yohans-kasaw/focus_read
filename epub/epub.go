package epub

import (
	"archive/zip"
	"fmt"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func New(file_path string) (*Epub, error) {
	epub := Epub{
		file_path: file_path,
	}

	epub.parseAndStore()
	return &epub, nil
}

func(e *Epub) parseAndStore(){
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

func(e *Epub) parseEpubToTexts(){
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
			e.Package.NavigationFile = item.Href
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
		e.Toc.FlatNavPoints = append(e.Toc.FlatNavPoints, nv)
		if len(nv.NavPoints) > 0 {
			e.parseFlatNavPoint(nv.NavPoints)
		}
	}
}

func extract_and_add(node *html.Node, arr *[]Text) {
	if node.Type == html.TextNode && node.Parent != nil {
		text := strings.TrimSpace(node.Data)
		if text != "" {
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
