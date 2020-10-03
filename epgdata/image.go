package epgdata

import (
	"encoding/xml"
	"regexp"
)

type Image struct {
	Source string
	Width  string
	Height string
}

func (i *Image) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	err := d.DecodeElement(&i.Source, &start)
	if err != nil {
		return err
	}

	i.parseImageSize()

	return nil
}

func (i *Image) parseImageSize() {
	re := regexp.MustCompile(`_(\d*)x(\d*).`)
	parts := re.FindStringSubmatch(i.Source)

	if len(parts) == 3 {
		i.Width = parts[1]
		i.Height = parts[2]
	}
}
