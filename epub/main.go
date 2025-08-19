package epub

import (
	"archive/zip"
)

type Epub struct {
	Rootfile struct {
		Path string `xml:"full-path,attr" json:"path"`
	} `xml:"rootfiles>rootfile"`

	FileMap map[string]*zip.File
}

type Opf struct {
	Version string `xml:"version,attr"`
}

func NewEpub(file_path string) (error, *Epub) {
	r, err := zip.OpenReader(file_path)
	if err != nil {
		return err, nil
	}
	defer r.Close()

	fileMap := make(map[string]*zip.File)
	for _, file := range r.File {
		fileMap[file.Name] = file
	}

	epub := Epub{
		FileMap: fileMap,
	}

	readXml(fileMap["META-INF/container.xml"], &epub)

	return nil, &epub
}
