package epgdata

import (
	"encoding/xml"
	"time"
)

type dateTime struct {
	time.Time
}

func (t *dateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02 15:04:05"
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}

	// correct missing time zone
	*t = dateTime{parse.Add(time.Hour * -2)}

	return nil
}

type tinyDate struct {
	time.Time
}

func (t *tinyDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "20060102"
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}

	*t = tinyDate{parse}

	return nil
}
