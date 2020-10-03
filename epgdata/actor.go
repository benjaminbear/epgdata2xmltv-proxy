package epgdata

import (
	"encoding/xml"
	"regexp"
	"strings"
)

type Actors struct {
	Actors []*Actor
}

type Actor struct {
	Actor string
	Role  string
}

func (a *Actors) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, " (voice)", "")
	splits := strings.Split(v, " - ")

	re := regexp.MustCompile(`^(.*)\s\((.*)\)`)

	for _, split := range splits {
		parts := re.FindStringSubmatch(split)

		if len(parts) > 0 {
			actor := &Actor{
				Actor: parts[1],
			}

			if len(parts) > 1 {
				actor.Role = parts[2]
			}

			a.Actors = append(a.Actors, actor)
		}
	}

	return nil
}
