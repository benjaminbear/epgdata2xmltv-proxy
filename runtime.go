package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/benjaminbear/epgdata2xmltv-proxy/config"
	"github.com/benjaminbear/epgdata2xmltv-proxy/epgdata"
	"github.com/benjaminbear/epgdata2xmltv-proxy/epgdownload"
	"github.com/benjaminbear/epgdata2xmltv-proxy/today"
	"github.com/benjaminbear/epgdata2xmltv-proxy/xmltv"
)

type RunTime struct {
	GenreMap    map[string]string
	CategoryMap map[string]string
	ChannelMap  map[string]string
	Channels    []*xmltv.Channel
	XMLTv       *xmltv.Tv
	EPGDays     []*epgdata.Pack
	Config      *config.Config
}

func (r *RunTime) initEPGDays() {
	r.EPGDays = make([]*epgdata.Pack, r.Config.Days)

	for i, _ := range r.EPGDays {
		r.EPGDays[i] = &epgdata.Pack{}
	}
}

func (r *RunTime) EPGCron() error {
	fmt.Println("Starting cron ... ")
	timeToday := today.New()

	// check day0, make day dance
	err := r.DayRotation(timeToday)
	if err != nil {
		return err
	}

	// download days
	err = epgdownload.DownloadEPG(r.Config.Pin, timeToday, r.Config.Days)
	if err != nil {
		return err
	}

	// parse newly downloaded days
	for i, epgDay := range r.EPGDays {
		if epgDay.Date == "" {
			matches, err := filepath.Glob(filepath.Join("epgdata_files", timeToday.GetDayPlus(i)+"_*_de_qy.xml"))
			if err != nil {
				return err
			}

			if len(matches) < 1 {
				return fmt.Errorf("epgfile for date %s not found", timeToday.GetDayPlus(i))
			}

			r.EPGDays[i], err = epgdata.ReadEPGDataFile(matches[0])
			if err != nil {
				return err
			}

			r.EPGDays[i].Date = timeToday.GetDayPlus(i)
		}
	}

	fmt.Println("finished!")

	return nil
}

func (r *RunTime) DayRotation(timeToday *today.Today) error {
	for i, epgDay := range r.EPGDays {
		// today found, make rotation
		if epgDay.Date == timeToday.String() {
			r.EPGDays = r.EPGDays[i:]

			for j := 0; j < i; j++ {
				r.EPGDays = append(r.EPGDays, &epgdata.Pack{})
			}

			return nil
		}

		// if epgDay.Date != today empty it
		r.EPGDays[i].Date = ""
	}

	// today not found
	return nil
}

func (r *RunTime) ParseChannels(channels *epgdata.Channels) {
	for _, chanData := range channels.ChannelsData {
		r.Channels = append(r.Channels, &xmltv.Channel{
			Id: chanData.TvChannelId,
			DisplayName: &xmltv.StdLangElement{
				Language: chanData.CountryDomain,
				Name:     chanData.TvChannelName,
			},
		})
	}
}

func (r *RunTime) Merge() error {
	// Create XMLTV File
	r.XMLTv = xmltv.NewXMLTVFile()
	r.XMLTv.Channels = r.Channels
	r.XMLTv.GeneratorInfoName = "epgdata2xmltv-proxy"
	r.XMLTv.GeneratorInfoURL = "https://github.com/benjaminbear/epgdata2xmltv-proxy"

	for _, epg := range r.EPGDays {
		for _, program := range epg.Programs {
			tvProgram := &xmltv.Program{
				Channel: program.TvChannelId,
				Start:   program.StartTime.In(r.Config.TimeZone).Format(xmltv.DateTimeFormat),
				Stop:    program.EndTime.In(r.Config.TimeZone).Format(xmltv.DateTimeFormat),
			}

			if program.Title != "" {
				tvProgram.Title = &xmltv.StdLangElement{
					Name:     program.Title,
					Language: r.ChannelMap[program.TvChannelId],
				}
			}

			if program.Subtitle != "" {
				tvProgram.SubTitle = &xmltv.StdLangElement{
					Name:     program.Subtitle,
					Language: r.ChannelMap[program.TvChannelId],
				}
			}

			if program.CommentLong != "" {
				tvProgram.Desc = &xmltv.StdLangElement{
					Name:     program.CommentLong,
					Language: r.ChannelMap[program.TvChannelId],
				}
			}

			if program.SeasonNum != 0 {
				tvProgram.EpisodeNum = &xmltv.EpisodeNum{
					Value:  fmt.Sprintf("%d.%d.", program.SeasonNum-1, program.EpisodeNum-1),
					System: "xmltv_ns",
				}
			}

			if program.ImageBig.Source != "" {
				tvProgram.Icon = &xmltv.Icon{
					Source: program.ImageBig.Source,
					Width:  program.ImageBig.Width,
					Height: program.ImageBig.Height,
				}
			}

			if program.Year != "" {
				tvProgram.Date = program.Year
			}

			if program.Country != "" {
				tvProgram.Country = &xmltv.StdLangElement{
					Name: program.Country,
				}
			}

			if program.TvShowLength != "" {
				tvProgram.Length = &xmltv.Length{
					Value: program.TvShowLength,
					Unit:  xmltv.Minutes,
				}
			}

			if program.AgeMarker != "" {
				tvProgram.Rating = append(tvProgram.Rating, &xmltv.Rating{
					Value:  program.AgeMarker,
					System: "FSK",
				})
			}

			program.ParseIncludeData(r.GenreMap, r.CategoryMap)
			if program.CategoryId != "" || program.GenreId != "" {
				if program.Category != "" {
					tvProgram.Categories = append(tvProgram.Categories, &xmltv.StdLangElement{
						Name:     program.Category,
						Language: r.ChannelMap[program.TvChannelId],
					})
				}

				if program.Genre != "" {
					tvProgram.Categories = append(tvProgram.Categories, &xmltv.StdLangElement{
						Name:     program.Genre,
						Language: r.ChannelMap[program.TvChannelId],
					})
				}
			}

			if program.HasCredits() {
				tvProgram.Credits = &xmltv.Credits{}

				if len(program.Moderator.List) > 0 {
					tvProgram.Credits.Presenters = program.Moderator.List
				}

				if len(program.StudioGuest.List) > 0 {
					tvProgram.Credits.Guests = program.StudioGuest.List
				}

				if len(program.Regisseurs.List) > 0 {
					tvProgram.Credits.Directors = program.Regisseurs.List
				}

				if len(program.Actors.Actors) > 0 {
					for _, actor := range program.Actors.Actors {
						a := &xmltv.Actor{
							Name: actor.Actor,
						}

						if actor.Role != "" {
							a.Role = actor.Role
						}

						tvProgram.Credits.Actors = append(tvProgram.Credits.Actors, a)
					}
				}
			}

			r.XMLTv.Programs = append(r.XMLTv.Programs, tvProgram)
		}
	}

	return nil
}

func (r *RunTime) XMLTvServer(w http.ResponseWriter, req *http.Request) {
	// Create XMLTv
	err := r.Merge()
	if err != nil {
		fmt.Println(err)
	}

	// Marshal
	err = r.XMLTv.WriteFile("latest.xml")
	if err != nil {
		fmt.Println(err)
	}

	http.ServeFile(w, req, "latest.xml")

	IPAddress := req.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = req.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = req.RemoteAddr
	}

	fmt.Println("Request from: ", IPAddress)
}
