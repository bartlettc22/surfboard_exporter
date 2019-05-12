package main

import (
  "net/http"
	"strings"

  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"

  // "github.com/davecgh/go-spew/spew"
)

// MetricsHandler is our http request handler
type MetricsHandler struct{}

func main() {

	var metricsHandler MetricsHandler
	http.Handle("/metrics", metricsHandler)

	log.Info("Listening on port 9040")
	http.ListenAndServe(":9040", nil)

	// spew.Dump("done")
}

// ServeHTTP scrapes the modem data and serves up the Prometheus metrics
func (m MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Info("metrics requested")
	// Parse the modem pages
	parseHelpPage()

	// Let promhttp serve up the metrics page
	promHandler := promhttp.Handler()
	promHandler.ServeHTTP(w, r)
}

// getPage downloads and returns the url document body
func getPage(url string) *goquery.Document {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("status code error: ", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// parseHelpPage parses the modem's help page for hardward/firmware info
func parseHelpPage() {

  // var labels prometheus.Labels
  labels := make(map[string]string)

	// Get the document
	page := getPage("http://192.168.100.1/cmHelpData.htm")

	// Get the tables on the page
	tables := page.Find("html body table")

	// The data we want is in the first table (before JS loads header), first row, first column
	html, _ := tables.Eq(0).Find("td").Eq(0).Html()

	// The element is a line-break separated string
	parts := strings.Split(html, "<br/>")

  // This will parse the fields with the info and get the text to the right
  // of the first colon (:), trimmed of whitespace
  getFieldValue := func(text string) string {
    return strings.TrimSpace(text[strings.Index(text, ":")+1 : len(text)])
  }

  // Create our label fields
  labels["model_name"] = getFieldValue(parts[0])
  labels["vendor_name"] = getFieldValue(parts[1])
  labels["firmware_name"] = getFieldValue(parts[2])
  labels["boot_version"] = getFieldValue(parts[3])
  labels["hardware_version"] = getFieldValue(parts[4])
  labels["serial_number"] = getFieldValue(parts[5])
  labels["firmware_build_time"] = getFieldValue(parts[6])

  // Set the metric
  modemInfo.With(labels).Set(1)
}
