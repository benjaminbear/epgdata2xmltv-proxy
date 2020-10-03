package epgdata

import (
	"encoding/xml"
	"io/ioutil"
)

type Categories struct {
	XMLName    xml.Name    `xml:"category"`
	Categories []*Category `xml:"data"`
}

type Category struct {
	CategoryId    string `xml:"ca0"`
	CategoryLong  string `xml:"ca1"`
	CategoryShort string `xml:"ca2"`
}

func NewCategories() *Categories {
	category := &Categories{}
	category.Categories = make([]*Category, 0)

	return category
}

func MarshalCategories(v interface{}) ([]byte, error) {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return data, err
	}

	data = append([]byte(xml.Header), data...)

	return data, err
}

func UnmarshalCategories(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func ReadCategoriesFile(path string) (categories *Categories, categoryMap map[string]string, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	categories = NewCategories()

	err = UnmarshalCategories(data, categories)
	if err != nil {
		return nil, nil, err
	}

	categoryMap = make(map[string]string)
	for _, category := range categories.Categories {
		categoryMap[category.CategoryId] = category.CategoryLong
	}

	return categories, categoryMap, nil
}

func WriteCategoriesFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
