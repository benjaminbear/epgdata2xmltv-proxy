package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Pin      string
	Days     int
	TimeZone *time.Location
	Crawler  bool
}

func ParseEnv() (conf *Config, err error) {
	conf = &Config{
		Pin:  os.Getenv("EPG2XMLTV_EPGDATA_PIN"),
		Days: 7,
	}

	if conf.Pin == "" {
		return conf, fmt.Errorf("no epgdata pin set")
	}

	dayStr := os.Getenv("EPG2XMLTV_DAYS")
	if dayStr != "" {
		conf.Days, err = strconv.Atoi(dayStr)
		if err != nil {
			return conf, err
		}
	}

	tz := os.Getenv("EPG2XMLTV_TIMEZONE")
	if tz == "" {
		conf.TimeZone = time.Now().Local().Location()
	} else {
		conf.TimeZone, err = time.LoadLocation(tz)
		if err != nil {
			return conf, err
		}
	}

	crw := os.Getenv("EPG2XMLTV_CRAWLER")
	if crw == "" {
		conf.Crawler = true
	} else {
		conf.Crawler, err = strconv.ParseBool(crw)
		if err != nil {
			return conf, err
		}
	}

	return conf, err
}
