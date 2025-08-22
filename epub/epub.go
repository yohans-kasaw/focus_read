package epub

import (
	"archive/zip"
	"fmt"
	"path"
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


	if err = epub.parseToc(); err != nil {
		return nil, err
	}

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
