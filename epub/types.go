package epub

import "archive/zip"

type Epub struct {
	ZipReader *zip.ReadCloser
	FileMap   map[string]*zip.File
	TextFiles []string

	Container Container
	Package   Package
	Toc       Toc
}

type Container struct {
	Rootfile struct {
		Path string `xml:"full-path,attr"`
	} `xml:"rootfiles>rootfile"`
}

type Package struct {
	NavigationFile string
	Title          string `xml:"metadata>title"`
	Version        string `xml:"version,attr"`
	Manifest       struct {
		Items []struct {
			Href       string `xml:"href,attr"`
			Id         string `xml:"id,attr"`
			Media_type string `xml:"media-type,attr"`
			// for version 3.0
			Properties string `xml:"properties"`
		} `xml:"item"`
	} `xml:"manifest"`
	Spine struct {
		Items []struct {
			Idref string `xml:"idref,attr"`
		} `xml:"itemref"`
	} `xml:"spine"`
}

type Toc struct {
	NavPoints []struct {
		NavLabel string `xml:"navLabel>text"`
		Content  struct {
			Src string `xml:"src,attr"`
		} `xml:"content"`
	} `xml:"navMap>navPoint"`
}
