package epub

import (
	"archive/zip"
	"fmt"
)

func New(file_path string) (*Epub, error) {
	r, err := zip.OpenReader(file_path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

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

	package_path, ok := e.FileMap[e.Container.Rootfile.Path]
	if !ok {
		return fmt.Errorf(
			"package file %s not found in EPUB",
			e.Container.Rootfile.Path,
		)
	}

	if err := readXml(package_path, &e.Package); err != nil {
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
