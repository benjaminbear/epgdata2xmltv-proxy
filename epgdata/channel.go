package epgdata

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
)

type Channels struct {
	XMLName      xml.Name       `xml:"channel"`
	ChannelsData []*ChannelData `xml:"data"`
}

type ChannelData struct {
	TvChannelName  string `xml:"ch0"`
	TvChannelShort string `xml:"ch1"`
	LanguageEn     string `xml:"ch2"`
	CountryDomain  string `xml:"ch3"`
	TvChannelId    string `xml:"ch4"`
	Sort           string `xml:"ch5"`
	PackageId      string `xml:"ch6"`
	Cni840f1       string `xml:"ch7"`
	CniVps         string `xml:"ch8"`
	Cni830f2       string `xml:"ch9"`
	Cnix26dw       string `xml:"ch10"`
	TvChannelDvb   string `xml:"ch11"`
	TvChannelType  string `xml:"ch12"`
}

func NewChannels() *Channels {
	channel := &Channels{}
	channel.ChannelsData = make([]*ChannelData, 0)

	return channel
}

func MarshalChannels(v interface{}) ([]byte, error) {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return data, err
	}

	data = append([]byte(xml.Header), data...)

	return data, err
}

func UnmarshalChannels(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func ReadChannelsFile() (channels *Channels, channelMap map[string]string, err error) {
	data, err := ioutil.ReadFile(filepath.Join(folderEPGInclude, fileEPGIncludeChannels))
	if err != nil {
		return nil, nil, err
	}

	channels = NewChannels()

	err = UnmarshalChannels(data, channels)
	if err != nil {
		return nil, nil, err
	}

	channelMap = make(map[string]string)
	for _, channel := range channels.ChannelsData {
		channelMap[channel.TvChannelId] = channel.CountryDomain
	}

	return channels, channelMap, nil
}

func WriteChannelsFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
