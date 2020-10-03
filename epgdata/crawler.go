package epgdata

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const imageBaseUrl = "https://cellular.images.dvbdata.com"

func (p *Program) GetAdditionalData() error {
	// Make HTTP request
	response, err := http.Get("https://m.hoerzu.de/tv-programm/" + p.BroadcastId + "/" + p.TvShowId + "/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("request returned", response.StatusCode)
	}

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	// Find and parse season number
	if p.EpisodeNum != 0 {
		selection := document.Find(".season-episode").First()
		splits := strings.Split(selection.Text(), ",")
		splits = strings.Split(splits[0], " ")

		if len(splits) > 1 {
			p.SeasonNum, err = strconv.Atoi(splits[1])
		}
	}

	// Find and print image URL
	if p.ImageBig.Source == "" {
		document.Find("meta").Each(func(i int, s *goquery.Selection) {
			op, _ := s.Attr("property")
			con, _ := s.Attr("content")
			if op == "og:image" {
				splits := strings.Split(con, "/")
				p.ImageBig.Source = con

				if len(splits) >= 3 && len(splits[3]) > 6 {
					p.ImageBig.Source = imageBaseUrl + "/" + splits[3] + "/" + splits[3] + "_320x240.jpg"
				}

				p.ImageBig.parseImageSize()
			}
		})
	}

	return nil
}
