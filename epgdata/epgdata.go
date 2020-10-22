package epgdata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	folderPersistence        = "persistence"
	filePersistence          = "epgdata.bin"
	folderEPGInclude         = "epgdata_includes"
	fileEPGIncludeGenres     = "genre.xml"
	fileEPGIncludeCategories = "category.xml"
	fileEPGIncludeChannels   = "channel_y.xml"
)

var (
	pathPersistence = filepath.Join(folderPersistence, filePersistence)
)

type Pack struct {
	XMLName  xml.Name   `xml:"pack"`
	Programs []*Program `xml:"data"`
	Date     string
}

type Program struct {
	BroadcastId       string          `xml:"d0"`
	TvShowId          string          `xml:"d1"`
	TvChannelId       string          `xml:"d2"`
	TvRegionId        string          `xml:"d3"`
	StartTime         dateTime        `xml:"d4"`
	EndTime           dateTime        `xml:"d5"`
	BroadcastDay      tinyDate        `xml:"d6"`
	TvShowLength      string          `xml:"d7"`
	Vps               string          `xml:"d8"`
	PrimeTime         string          `xml:"d9"`
	CategoryId        string          `xml:"d10"`
	TechnicsBw        string          `xml:"d11"`
	TechnicsCoChannel string          `xml:"d12"`
	TechnicsVt150     string          `xml:"d13"`
	TechnicsCoded     string          `xml:"d14"`
	TechnicsBlind     string          `xml:"d15"`
	AgeMarker         string          `xml:"d16"`
	LiveId            string          `xml:"d17"`
	TipFlag           string          `xml:"d18"`
	Title             string          `xml:"d19"`
	Subtitle          string          `xml:"d20"`
	CommentLong       string          `xml:"d21"`
	CommentMiddle     string          `xml:"d22"`
	CommentShort      string          `xml:"d23"`
	Themes            string          `xml:"d24"`
	GenreId           string          `xml:"d25"`
	EpisodeNum        int             `xml:"d26"`
	TechnicsStereo    string          `xml:"d27"`
	TechnicsDolby     string          `xml:"d28"`
	TechnicsWide      string          `xml:"d29"`
	TvdTotalValue     string          `xml:"d30"`
	Attribute         string          `xml:"d31"`
	Country           string          `xml:"d32"`
	Year              string          `xml:"d33"`
	Moderator         *StdListElement `xml:"d34"`
	StudioGuest       *StdListElement `xml:"d35"`
	Regisseurs        *StdListElement `xml:"d36"`
	Actors            *Actors         `xml:"d37"`
	ImageSmall        *Image          `xml:"d38"`
	ImageMiddle       *Image          `xml:"d39"`
	ImageBig          *Image          `xml:"d40"`
	SeasonNum         int
	Category          string
	Genre             string
}

func (p *Program) HasCredits() bool {
	if len(p.Moderator.List) > 0 || len(p.StudioGuest.List) > 0 ||
		len(p.Regisseurs.List) > 0 || len(p.Actors.Actors) > 0 {
		return true
	}

	return false
}

func (p *Program) ParseIncludeData(genres map[string]string, categories map[string]string) {
	p.Genre = genres[p.GenreId]
	p.Category = categories[p.CategoryId]
}

func NewEPGData() *Pack {
	epgData := &Pack{}
	epgData.Programs = make([]*Program, 0)

	return epgData
}

func UnmarshalEPGData(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func ReadEPGDataFile(path string, crawler bool) (*Pack, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	epgData := NewEPGData()

	err = UnmarshalEPGData(data, epgData)
	if err != nil {
		return nil, err
	}

	if !crawler {
		return epgData, nil
	}

	sem := make(chan int, 10)

	for _, program := range epgData.Programs {
		if program.EpisodeNum == 0 && program.ImageBig.Source != "" {
			continue
		}

		sem <- 1

		go func(p *Program) {
			err := p.GetAdditionalData()
			if err != nil {
				fmt.Println(err)
			}

			<-sem
		}(program)
	}

	return epgData, nil
}

func Save(data interface{}) (err error) {
	if _, err = os.Stat(pathPersistence); !os.IsNotExist(err) {
		err = os.Remove(pathPersistence)
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(folderPersistence); os.IsNotExist(err) {
		err = os.MkdirAll(folderPersistence, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(pathPersistence)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(b))

	return err
}

func Load(packs []*Pack) (err error) {
	if _, err = os.Stat(pathPersistence); os.IsNotExist(err) {
		return nil
	}

	fmt.Println("Loading persistent data from disk")

	f, err := os.Open(pathPersistence)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(bufio.NewReader(f)).Decode(&packs)

	return err
}
