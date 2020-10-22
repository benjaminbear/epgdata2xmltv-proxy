# EPGDATA to XMLTV Proxy Server

![Travis build status](https://travis-ci.com/benjaminbear/epgdata2xmltv-proxy.svg?branch=master)
![Docker build status](https://img.shields.io/docker/cloud/build/bbaerthlein/epgdata2xmltv-proxy)
![Docker build automated](https://img.shields.io/docker/cloud/automated/bbaerthlein/epgdata2xmltv-proxy)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/benjaminbear/epgdata2xmltv-proxy)
![Go version](https://img.shields.io/github/go-mod/go-version/benjaminbear/epgdata2xmltv-proxy?filename=go.mod)
![License](https://img.shields.io/github/license/benjaminbear/epgdata2xmltv-proxy)

This docker service grabs **epgdata** from a source (hoerzu), adds further information from the web and deploys **xmltv** format via http server.

## Usage

Pull repository

```bash
docker pull bbaerthlein/epgdata2xmltv-proxy
```


Run container:

```bash
docker run -p 8080:8080 bbaerthlein/epgdata2xmltv-proxy
```

## Environment variables

```
EPG2XMLTV_EPGDATA_PIN: (mandatory) Your epgdata pin, mandatory
EPG2XMLTV_DAYS: (optional) Count of days to download from your epg source and serve: 1-7
EPG2XMLTV_TIMEZONE: (optional) Set your timezone for correct airtime: e.g. "Europe/Berlin"
EPG2XMLTV_CRAWLER: (optional) Enable/disable webcrawler (much faster on start)
```

## Mount directories

If you want persist downloaded epgdata files add:

```
-v /your/local/path:/epg/epgdata_files
```

to persist working data of the proxy (absolutely recommended):

```
-v /your/local/path:/epg/persistence
```

to add your own genre.xml, category.xml and channel_y.xml definition files:

```
-v /your/local/definitions:/epg/epgdata_includes
```

## Docker hub

https://hub.docker.com/r/bbaerthlein/epgdata2xmltv-proxy