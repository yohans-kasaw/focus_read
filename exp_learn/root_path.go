package explearn

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
)

const ContainerXMLFile = "META-INF/container.xml"

type Container struct {
	Rootfiles []struct {
		FullPath string `xml:"full-path,attr"`
	} `xml:"rootfiles>rootfile"`
}

func Root_path() {
	// get absolute file path
	abs_file_path := "/home/yohansh/focus_read/test_file/test.epub"

	r, error := zip.OpenReader(abs_file_path)
	if error != nil {
		fmt.Println(error)
	}
	defer r.Close()

	for _, file := range r.File {
		if file.Name != ContainerXMLFile {
			continue
		}

		f_reader, error := file.Open()

		if error != nil {
			fmt.Println(error)
		}

		xml_byte, error := io.ReadAll(f_reader)

		if error != nil {
			fmt.Println(error)
		}

		var container Container
		err := xml.Unmarshal(xml_byte, &container)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(container.Rootfiles[0].FullPath)

	}

	// unzip it
	// extract and cache it

	// open the container.xml
	// parse the xml
}
