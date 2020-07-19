package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"regexp"
)

// GetPage downloads and returns the url document body
func GetPage(url url.URL) (*goquery.Document, error) {

	res, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("Error fetching URL %s: %v", url.String(), err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("URL returned status code %s: %s", res.Status, url.String())
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing document %s: %v", url.String(), err)
	}

	return doc, nil
}

func GetFieldValue(text string) string {

	reg, err := regexp.Compile("[0-9A-Za-z.]+")
	if err != nil {
		log.Warn(err)
		return ""
	}
	text = reg.FindString(text)

	return text
}
