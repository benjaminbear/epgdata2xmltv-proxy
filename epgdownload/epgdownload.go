package epgdownload

import (
	"context"
	"crypto/tls"
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

const (
	epgURL        = "http://www.epgdata.com/index.php"
	FolderEPGData = "epgdata_files"
)

type EPGDownloader struct {
	Pin         string
	TimeToday   *today.Today
	Days        int
	InsecureTLS bool
}

func (e *EPGDownloader) DownloadEPG() error {
	err := e.removeDeprecated()
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "epgdataproxy")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	for i := 0; i < e.Days; i++ {
		matches, err := filepath.Glob(filepath.Join(FolderEPGData, e.TimeToday.GetDayPlus(i)+"_*_de_qy.xml"))
		if err != nil {
			return err
		}

		if len(matches) > 0 {
			fmt.Println("EPG File", e.TimeToday.GetDayPlus(i), "already downloaded, skipping download")
			continue
		}

		err = e.downloadFile(i, dir)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully downloaded epg for day %s\n", e.TimeToday.GetDayPlus(i))

		ext, err := fastzip.NewExtractor(filepath.Join(dir, fmt.Sprintf("%d.zip", i)), FolderEPGData)
		if err != nil {
			return err
		}
		defer ext.Close()

		if err = ext.Extract(context.Background()); err != nil {
			return err
		}

		fmt.Printf("Successfully extracted epg for day %s\n", e.TimeToday.GetDayPlus(i))
	}

	return nil
}

func (e *EPGDownloader) downloadFile(dayPlus int, tempDir string) error {
	req, err := http.NewRequest(http.MethodGet, epgURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("action", "sendPackage")
	q.Add("iOEM", "vdr")
	q.Add("dayOffset", fmt.Sprintf("%d", dayPlus))
	q.Add("pin", e.Pin)
	q.Add("dataType", "xml")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	//Get the response bytes from the url
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: e.InsecureTLS},
	}

	client := &http.Client{
		Transport: tr,
	}

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

func (e *EPGDownloader) removeDeprecated() error {
	matches, err := filepath.Glob(filepath.Join(FolderEPGData, "*_*_de_qy.xml"))
	if err != nil {
		return err
	}

	for _, file := range matches {
		fileBase := filepath.Base(file)
		fileDate, err := strconv.Atoi(strings.Split(fileBase, "_")[0])
		if err != nil {
			return err
		}

		dateToday, err := e.TimeToday.GetInt()
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
