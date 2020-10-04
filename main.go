package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/benjaminbear/epgdata2xmltv-proxy/config"

	"github.com/robfig/cron/v3"

	"github.com/benjaminbear/epgdata2xmltv-proxy/epgdata"
)

var Version = "undefined"

func main() {
	fmt.Println("Version:", Version)

	// Parse config from environment
	conf, err := config.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	runtime := &RunTime{
		Config: conf,
	}

	runtime.initEPGDays()

	var channels *epgdata.Channels
	_, runtime.GenreMap, err = epgdata.ReadGenresFile(filepath.Join("epgdata_includes", "genre.xml"))
	if err != nil {
		log.Fatal(err)
	}

	_, runtime.CategoryMap, err = epgdata.ReadCategoriesFile(filepath.Join("epgdata_includes", "category.xml"))
	if err != nil {
		log.Fatal(err)
	}

	channels, runtime.ChannelMap, err = epgdata.ReadChannelsFile(filepath.Join("epgdata_includes", "channel_y.xml"))
	if err != nil {
		log.Fatal(err)
	}

	runtime.ParseChannels(channels)

	// Start cron once
	err = runtime.EPGCron()
	if err != nil {
		fmt.Println(err)
	}

	// Add cron process
	c := cron.New(cron.WithLocation(runtime.Config.TimeZone))
	c.AddFunc("42 3 * * *", func() {
		err := runtime.EPGCron()
		if err != nil {
			fmt.Println(err)
		}
	})
	c.Start()

	// Run Server
	http.HandleFunc("/epg", runtime.XMLTvServer)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
