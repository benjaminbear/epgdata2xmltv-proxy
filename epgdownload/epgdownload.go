package epgdownload

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/benjaminbear/epgdata2xmltv-proxy/today"

	"github.com/saracen/fastzip"
)

const EGPURL = "http://www.epgdata.com/index.php"

type epgDownloader struct {
	pin       string
	timeToday *today.Today
	days      int
}

func DownloadEPG(pin string, timeToday *today.Today, days int) error {
	e := epgDownloader{
		pin:       pin,
		timeToday: timeToday,
		days:      days,
	}

	err := e.removeDeprecated()
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "epgdataproxy")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	for i := 0; i < e.days; i++ {
		matches, err := filepath.Glob(filepath.Join("epgdata_files", timeToday.GetDayPlus(i)+"_*_de_qy.xml"))
		if err != nil {
			return err
		}

		if len(matches) > 0 {
			fmt.Println("EPG File", timeToday.GetDayPlus(i), "already downloaded, skipping download")
			continue
		}

		err = e.downloadFile(i, dir)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully downloaded epg for day %s\n", timeToday.GetDayPlus(i))

		e, err := fastzip.NewExtractor(filepath.Join(dir, fmt.Sprintf("%d.zip", i)), "epgdata_files")
		if err != nil {
			return err
		}
		defer e.Close()

		if err = e.Extract(); err != nil {
			return err
		}

		fmt.Printf("Successfully extracted epg for day %s\n", timeToday.GetDayPlus(i))
	}

	return nil
}

func (e *epgDownloader) downloadFile(dayPlus int, tempDir string) error {
	req, err := http.NewRequest(http.MethodGet, EGPURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("action", "sendPackage")
	q.Add("iOEM", "vdr")
	q.Add("dayOffset", fmt.Sprintf("%d", dayPlus))
	q.Add("pin", e.pin)
	q.Add("dataType", "xml")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	//Get the response bytes from the url
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("request returned response %v", response.StatusCode)
	}

	//Create a empty file
	file, err := os.Create(filepath.Join(tempDir, fmt.Sprintf("%d.zip", dayPlus)))
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (e *epgDownloader) removeDeprecated() error {
	matches, err := filepath.Glob(filepath.Join("epgdata_files", "*_*_de_qy.xml"))
	if err != nil {
		return err
	}

	for _, file := range matches {
		fileBase := filepath.Base(file)
		fileDate, err := strconv.Atoi(strings.Split(fileBase, "_")[0])
		if err != nil {
			return err
		}

		dateToday, err := e.timeToday.GetInt()
		if err != nil {
			return err
		}

		if fileDate < dateToday {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
