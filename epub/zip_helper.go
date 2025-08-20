package epub

import (
	"archive/zip"
	"encoding/xml"
)

func readXml(file *zip.File, v interface{}) error {
	f_reader, err := file.Open()
	if err != nil {
		return err
	}
	defer f_reader.Close()

	return xml.NewDecoder(f_reader).Decode(v)
}
