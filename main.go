package main

import (
  "net/http"
  "regexp"
  "strconv"
	"strings"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"

  "github.com/davecgh/go-spew/spew"
)

// MetricsHandler is our http request handler
type MetricsHandler struct{}

type downstream struct {
  channelID string
  frequency float64
  snr float64
  modulation string
  power float64
}

type upstream struct {
  channelID string
  frequency float64
  rangingServiceID string
  symbolRate float64
  power float64
  modulation string
  rangingStatus float64
}

func main() {

  // parseSignalPage()
	// parseHelpPage()

	var metricsHandler MetricsHandler
	http.Handle("/metrics", metricsHandler)

	log.Info("Listening on port 9040")
	http.ListenAndServe(":9040", nil)

	spew.Dump("done")
}

// ServeHTTP scrapes the modem data and serves up the Prometheus metrics
func (m MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Parse the modem pages
  parseSignalPage()
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

// parseSignalPage parses the modem's signal page for upstream/downstream info
func parseSignalPage() {

	// Get the document
	page := getPage("http://192.168.100.1/cmSignalData.htm")

	// Get the centers on the page (these contain the tables)
	centers := page.Find("html body").ChildrenFiltered("center")

  parseDownstream(centers.Eq(0).ChildrenFiltered("table"))
  parseUpstream(centers.Eq(1).ChildrenFiltered("table"))

}

func parseDownstream(s *goquery.Selection) {

  // Downstream data is in the first table
  // There are 4 downstream channels to parse. We'll do them at the same time
  var ds [4]downstream

  channels := s.Find("tr").Eq(1).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := getFieldValue(channels.Eq(1+i).Text())
    ds[i].channelID = value
  }

  frequency := s.Find("tr").Eq(2).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(frequency.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }

    ds[i].frequency = value
  }

  snr := s.Find("tr").Eq(3).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(snr.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }

    ds[i].snr = value
  }

  mod := s.Find("tr").Eq(4).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := getFieldValue(mod.Eq(1+i).Text())
    ds[i].modulation = value
  }

  power := s.Find("tr").Eq(5).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(power.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }

    ds[i].power = value
  }

  for _, d := range ds {
    modemDownstreamFrequency.With(prometheus.Labels{"channel_id":d.channelID}).Set(d.frequency)
    modemDownstreamSNR.With(prometheus.Labels{"channel_id":d.channelID}).Set(d.snr)
    modemDownstreamPower.With(prometheus.Labels{"channel_id":d.channelID}).Set(d.power)
    modemDownstreamModulation.With(prometheus.Labels{"channel_id":d.channelID,"modulation_method":d.modulation}).Set(1)
  }
}

func parseUpstream(s *goquery.Selection) {

  // Downstream data is in the first table
  // There are 4 downstream channels to parse. We'll do them at the same time

  // type upstream struct {
  //   channelID string
  //   frequency float64
  //   rangingServiceID string
  //   symbolRate float64
  //   power float64
  //   modulation string
  //   rangingStatus int
  // }

  var us [4]upstream

  channels := s.Find("tr").Eq(1).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := getFieldValue(channels.Eq(1+i).Text())
    us[i].channelID = value
  }

  frequency := s.Find("tr").Eq(2).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(frequency.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }
    us[i].frequency = value
  }

  rsid := s.Find("tr").Eq(3).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := getFieldValue(rsid.Eq(1+i).Text())
    us[i].rangingServiceID = value
  }

  sr := s.Find("tr").Eq(4).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(sr.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }
    us[i].symbolRate = value
  }

  power := s.Find("tr").Eq(5).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value, err := strconv.ParseFloat(getFieldValue(power.Eq(1+i).Text()), 64)
    if err != nil {
      log.Warn(err)
    }
    us[i].power = value
  }

  mod := s.Find("tr").Eq(6).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := mod.Eq(1+i).Text()
    value = strings.ReplaceAll(value, "<br/>", "")
    value = strings.ReplaceAll(value, "\n", ",")
    value = strings.Trim(value, ",\u00a0")
    us[i].modulation = value
  }

  status := s.Find("tr").Eq(7).ChildrenFiltered("td")
  for i := 0; i < 4; i++ {
    value := status.Eq(1+i).Text()
    value = strings.Trim(value, "\u00a0")
    if value == "Success" {
      us[i].rangingStatus = 1
    } else {
      us[i].rangingStatus = 0
    }
  }

  for _, u := range us {
    modemUpstreamChannelInfo.With(prometheus.Labels{"channel_id":u.channelID,"ranging_service_id":u.rangingServiceID,"modulation_methods":u.modulation}).Set(1)
    modemUpstreamFrequency.With(prometheus.Labels{"channel_id":u.channelID}).Set(u.frequency)
    modemUpstreamSymbolRate.With(prometheus.Labels{"channel_id":u.channelID}).Set(u.symbolRate)
    modemUpstreamPower.With(prometheus.Labels{"channel_id":u.channelID}).Set(u.power)
    modemUpstreamRangingStatus.With(prometheus.Labels{"channel_id":u.channelID}).Set(u.rangingStatus)
  }

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
  getInfoValue := func(text string) string {
    return strings.TrimSpace(text[strings.Index(text, ":")+1 : len(text)])
  }

  // Create our label fields
  labels["model_name"] = getInfoValue(parts[0])
  labels["vendor_name"] = getInfoValue(parts[1])
  labels["firmware_name"] = getInfoValue(parts[2])
  labels["boot_version"] = getInfoValue(parts[3])
  labels["hardware_version"] = getInfoValue(parts[4])
  labels["serial_number"] = getInfoValue(parts[5])
  labels["firmware_build_time"] = getInfoValue(parts[6])

  // Set the metric
  modemInfo.With(labels).Set(1)
}

func getFieldValue(text string) string {

  reg, err := regexp.Compile("[0-9A-Za-z.]+")
  if err != nil {
      log.Warn(err)
      return ""
  }
  text = reg.FindString(text)

  return text
}
