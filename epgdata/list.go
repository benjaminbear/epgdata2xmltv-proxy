package epgdata

import (
	"encoding/xml"
	"strings"
)

type StdListElement struct {
	List []string
}

func (r *StdListElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	if v == "" {
		return nil
	}

	splits := strings.Split(v, "|")
	r.List = append(r.List, splits...)

	return nil
}
