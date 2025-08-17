package explearn

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"path"
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

func readXml(file *zip.File, v any) error {
	f_reader, err := file.Open()
	if err != nil {
		return err
	}
	defer f_reader.Close()

	return xml.NewDecoder(f_reader).Decode(v)
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

	container_file := fileMap["META-INF/container.xml"]
	readXml(container_file, &epub)

	return nil, &epub
}

func Main_func() {
	abs_file_path := "./test_file/test.epub"

	err, epub := NewEpub(abs_file_path)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(epub.Rootfile.Path)
	root_dir := path.Dir(epub.Rootfile.Path)
	fmt.Println(root_dir)
}
