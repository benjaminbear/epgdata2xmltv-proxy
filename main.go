package main

import (
	"fmt"
	"log"
	"net/http"

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

	_, runtime.GenreMap, err = epgdata.ReadGenresFile()
	if err != nil {
		log.Fatal(err)
	}

	_, runtime.CategoryMap, err = epgdata.ReadCategoriesFile()
	if err != nil {
		log.Fatal(err)
	}

	var channels *epgdata.Channels
	channels, runtime.ChannelMap, err = epgdata.ReadChannelsFile()
	if err != nil {
		log.Fatal(err)
	}

	runtime.ParseChannels(channels)

	// Load persistence if there
	err = epgdata.Load(runtime.EPGDays)
	if err != nil {
		log.Fatal(err)
	}

	// Start cron once
	err = runtime.EPGCron()
	if err != nil {
		fmt.Println(err)
	}

	// Add cron process
	c := cron.New(cron.WithLocation(runtime.Config.TimeZone))
	c.AddFunc("12 6 * * *", func() {
		err := runtime.EPGCron()
		if err != nil {
			fmt.Println(err)
		}
	})
	c.Start()

	// Run Server
	http.HandleFunc("/epg", runtime.XMLTvServer)
	fmt.Println("Webserver ready.")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
