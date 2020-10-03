package xmltv

import (
	"encoding/xml"
	"io/ioutil"
)

// Scheme
// https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd

const DateTimeFormat = "20060102150405 -0700"

type Tv struct {
	XMLName           xml.Name   `xml:"tv"`
	Channels          []*Channel `xml:"channel"`
	Programs          []*Program `xml:"programme"`
	Date              string     `xml:"date,attr,omitempty"`
	SourceInfoUrl     string     `xml:"source-info-url,attr,omitempty"`
	SourceInfoName    string     `xml:"source-info-name,attr,omitempty"`
	SourceDataURL     string     `xml:"source-data-url,attr,omitempty"`
	GeneratorInfoName string     `xml:"generator-info-name,attr,omitempty"`
	GeneratorInfoURL  string     `xml:"generator-info-url,attr,omitempty"`
}

type Channel struct {
	DisplayName *StdLangElement `xml:"display-name"`
	Icon        *Icon           `xml:"icon,omitempty"`
	URL         string          `xml:"url,omitempty"`
	Id          string          `xml:"id,attr"`
}

type StdLangElement struct {
	Name     string `xml:",chardata"`
	Language string `xml:"lang,attr,omitempty"`
}

type Program struct {
	Title           *StdLangElement   `xml:"title"`
	SubTitle        *StdLangElement   `xml:"sub-title,omitempty"`
	Desc            *StdLangElement   `xml:"desc,omitempty"`
	Credits         *Credits          `xml:"credits,omitempty"`
	Date            string            `xml:"date,omitempty"`
	Categories      []*StdLangElement `xml:"category,omitempty"`
	Keywords        []*StdLangElement `xml:"keyword,omitempty"`
	Language        *StdLangElement   `xml:"language,omitempty"`
	OrigLanguage    *StdLangElement   `xml:"orig-language,omitempty"`
	Length          *Length           `xml:"length,omitempty"`
	Icon            *Icon             `xml:"icon,omitempty"`
	URL             string            `xml:"url,omitempty"`
	Country         *StdLangElement   `xml:"country,omitempty"`
	EpisodeNum      *EpisodeNum       `xml:"episode-num,omitempty"`
	Video           *Video            `xml:"video,omitempty"`
	Audio           *Audio            `xml:"audio,omitempty"`
	PreviouslyShown string            `xml:"previously-shown,omitempty"`
	Premiere        *StdLangElement   `xml:"premiere,omitempty"`
	LastChance      *StdLangElement   `xml:"last-chance,omitempty"`
	New             string            `xml:"new,omitempty"`
	Subtitles       []*Subtitles      `xml:"subtitles,omitempty"`
	Rating          []*Rating         `xml:"rating,omitempty"`
	StarRating      []*Rating         `xml:"star-rating,omitempty"`
	Review          []*Review         `xml:"review,omitempty"`
	Start           string            `xml:"start,attr"`
	Stop            string            `xml:"stop,attr,omitempty"`
	PdcStart        string            `xml:"pdc-start,attr,omitempty"`
	VpsStart        string            `xml:"vps-start,attr,omitempty"`
	ShowView        string            `xml:"showview,attr,omitempty"`
	VideoPlus       string            `xml:"videoplus,attr,omitempty"`
	Channel         string            `xml:"channel,attr"`
	ClumpIdx        string            `xml:"clumpidx,attr,omitempty"`
}

type Credits struct {
	Directors    []string `xml:"director,omitempty"`
	Actors       []*Actor `xml:"actor,omitempty"`
	Writers      []string `xml:"writer,omitempty"`
	Adapters     []string `xml:"adapter,omitempty"`
	Producers    []string `xml:"producer,omitempty"`
	Composers    []string `xml:"composer,omitempty"`
	Editors      []string `xml:"editor,omitempty"`
	Presenters   []string `xml:"presenter,omitempty"`
	Commentators []string `xml:"commentator,omitempty"`
	Guests       []string `xml:"guest,omitempty"`
}

type Actor struct {
	Name string `xml:",chardata"`
	Role string `xml:"role,attr,omitempty"`
}

type Length struct {
	Value string   `xml:",chardata"`
	Unit  TimeUnit `xml:"units,attr"`
}

type Icon struct {
	Text   string `xml:",chardata"`
	Source string `xml:"src,attr"`
	Width  string `xml:"width,attr,omitempty"`
	Height string `xml:"height,attr,omitempty"`
}

type EpisodeNum struct {
	Value  string `xml:",chardata"`
	System string `xml:"system,attr"`
}

type Video struct {
	Present string `xml:"present,omitempty"`
	Colour  string `xml:"colour,omitempty"`
	Aspect  string `xml:"aspect,omitempty"`
	Quality string `xml:"quality,omitempty"`
}

type Audio struct {
	Present string `xml:"present,omitempty"`
	Stereo  string `xml:"stereo,omitempty"`
}

type PreviouslyShown struct {
	Start   string `xml:"start,attr,omitempty"`
	Channel string `xml:"channel,attr,omitempty"`
}

type Subtitles struct {
	Language *StdLangElement `xml:"language,omitempty"`
	Type     SubtitleType    `xml:"type,omitempty"`
}

type Rating struct {
	Value  string `xml:"value"`
	Icon   *Icon  `xml:"icon,omitempty"`
	System string `xml:"system,attr,omitempty"`
}

type Review struct {
	Text     string
	Type     ReviewType `xml:"type,attr"`
	Source   string     `xml:"source,attr,omitempty"`
	Reviewer string     `xml:"reviewer,attr,omitempty"`
	Lang     string     `xml:"lang,attr,omitempty"`
}

type TimeUnit string

const (
	Seconds TimeUnit = "seconds"
	Minutes TimeUnit = "minutes"
	Hours   TimeUnit = "hours"
)

type SubtitleType string

const (
	TeleText   SubtitleType = "teletext"
	OnScreen   SubtitleType = "onscreen"
	DeafSigned SubtitleType = "deaf-signed"
)

type ReviewType string

const (
	Text ReviewType = "text"
	Url  ReviewType = "url"
)

func NewXMLTVFile() *Tv {
	xmlTv := &Tv{}
	xmlTv.Channels = make([]*Channel, 0)
	xmlTv.Programs = make([]*Program, 0)

	return xmlTv
}

func (tv *Tv) Marshal() ([]byte, error) {
	data, err := xml.MarshalIndent(tv, "", "  ")
	if err != nil {
		return data, err
	}

	data = append([]byte(xml.Header), data...)

	return data, err
}

func Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func ReadFile(path string) (*Tv, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	xmltvf := NewXMLTVFile()

	err = Unmarshal(data, xmltvf)
	if err != nil {
		return nil, err
	}

	return xmltvf, nil
}

func (tv *Tv) WriteFile(path string) error {
	data, err := tv.Marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
