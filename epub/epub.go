package epub

import (
	"archive/zip"
	"fmt"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func New(file_path string) (*Epub, error) {
	r, err := zip.OpenReader(file_path)
	if err != nil {
		return nil, err
	}

	epub := Epub{
		ZipReader: r,
		FileMap:   make(map[string]*zip.File),
	}

	for _, file := range r.File {
		epub.FileMap[file.Name] = file
	}

	if err = epub.parseContainer(); err != nil {
		return nil, err
	}

	if err = epub.parsePackage(); err != nil {
		return nil, err
	}

	if err = epub.parseToc(); err != nil {
		return nil, err
	}

	epub.parseFlatNavPoint(epub.Toc.NavPoints)

	epub.constructTextFiles()

	return &epub, nil
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

func cli_show(file *zip.File) {
	r, _ := file.Open()
	z := html.NewTokenizer(r)
	for {
		token_type := z.Next()
		if token_type == html.ErrorToken {
			break
		}

		if token_type == html.TextToken {
			text := z.Token().Data
			text = strings.TrimSpace(text)
			if text != "" {
				fmt.Println(text)
			}
		}
	}
}

func (e *Epub) HtmlLineByLine() {
	for _, points := range e.Toc.FlatNavPoints[3:5] {
		file, ok := e.FileMap[points.Content.Src]
		if !ok {
			fmt.Println("hurray not found")
			fmt.Println("this file not found", points.NavLabel)
		} else {
			fmt.Println("---------------")
			fmt.Println(points.NavLabel)
			fmt.Println(points.Content.Src)
			fmt.Println("---------------")
			cli_show(file)
		}

	}
}
