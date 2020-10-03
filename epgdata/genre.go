package epgdata

import (
	"encoding/xml"
	"io/ioutil"
)

type Genres struct {
	XMLName xml.Name `xml:"genre"`
	Genres  []*Genre `xml:"data"`
}

type Genre struct {
	GenreId   string `xml:"g0"`
	GenreName string `xml:"g1"`
}

func NewGenres() *Genres {
	genre := &Genres{}
	genre.Genres = make([]*Genre, 0)

	return genre
}

func MarshalGenres(v interface{}) ([]byte, error) {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return data, err
	}

	data = append([]byte(xml.Header), data...)

	return data, err
}

func UnmarshalGenres(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func ReadGenresFile(path string) (genres *Genres, genreMap map[string]string, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	genres = NewGenres()
	err = UnmarshalGenres(data, genres)
	if err != nil {
		return nil, nil, err
	}

	genreMap = make(map[string]string)
	for _, genre := range genres.Genres {
		genreMap[genre.GenreId] = genre.GenreName
	}

	return genres, genreMap, nil
}

func WriteGenresFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
