package epub

import (
	"fmt"
	"path"
	"strings"
	"archive/zip"
	"golang.org/x/net/html"
)

func New(filePath string) (*Epub, error) {
	epub := &Epub{
		filePath: filePath,
	}

	if err := epub.parse(); err != nil {
		return nil, err
	}

	return epub, nil
}

func (e *Epub) parse() error {
	r, err := zip.OpenReader(e.filePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	e.ZipReader = r
	defer r.Close()

	e.FileMap = make(map[string]*zip.File)
	for _, file := range r.File {
		e.FileMap[file.Name] = file
	}

	if err := e.parseContainer(); err != nil {
		return fmt.Errorf("failed to parse container: %w", err)
	}

	if err := e.parsePackage(); err != nil {
		return fmt.Errorf("failed to parse package: %w", err)
	}

	if err := e.parseToc(); err != nil {
		return fmt.Errorf("failed to parse table of contents: %w", err)
	}

	e.parseFlatNavPoint(e.Toc.NavPoints)

	e.constructTextFiles()

	if err := e.parseEpubToTexts(); err != nil {
		return fmt.Errorf("failed to parse texts: %w", err)
	}

	return nil
}

func (e *Epub) parseEpubToTexts() error {
	var arr []Text

	for _, points := range e.Toc.FlatNavPoints {
		file, ok := e.FileMap[points.Content.Src]
		if !ok {
			continue
		}

		r, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file.Name, err)
		}
		defer r.Close()

		doc, err := html.Parse(r)
		if err != nil {
			return fmt.Errorf("failed to parse HTML for %s: %w", file.Name, err)
		}

		extractAndAdd(doc, &arr)
	}

	e.Texts = arr
	return nil
}

func (e *Epub) parseContainer() error {
	filePath, ok := e.FileMap["META-INF/container.xml"]
	if !ok {
		return fmt.Errorf("META-INF/container.xml not found in EPUB")
	}
	return readXml(filePath, &e.Container)
}

func (e *Epub) parsePackage() error {
	if e.Container.Rootfile.Path == "" {
		return fmt.Errorf("rootfile path not found in container.xml")
	}

	packageFile, ok := e.FileMap[e.Container.Rootfile.Path]
	if !ok {
		return fmt.Errorf("package file %s not found in EPUB", e.Container.Rootfile.Path)
	}

	if err := readXml(packageFile, &e.Package); err != nil {
		return fmt.Errorf("failed to read package XML: %w", err)
	}

	for _, item := range e.Package.Manifest.Items {
		if item.Media_type == "application/x-dtbncx+xml" || item.Properties == "nav" {
			navPath := path.Join(path.Dir(e.Container.Rootfile.Path), item.Href)
			e.Package.NavigationFile = navPath
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
		return fmt.Errorf("NavigationFile is null")
	}

	navZipFile, ok := e.FileMap[e.Package.NavigationFile]
	if !ok {
		return fmt.Errorf("nav file %s not found in EPUB", e.Container.Rootfile.Path)
	}

	return readXml(navZipFile, &e.Toc)
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

func extractAndAdd(node *html.Node, arr *[]Text) {
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
		extractAndAdd(c, arr)
	}
}

func (e *Epub) getFullPath(s string) string {
	filePath, _, _ := strings.Cut(s, "#")
	fullPath := path.Join(path.Dir(e.Container.Rootfile.Path), filePath)
	return fullPath
}
