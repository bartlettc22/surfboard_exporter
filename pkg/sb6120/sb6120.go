package sb6120

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bartlettc22/surfboard_exporter/pkg/modem"
	"github.com/bartlettc22/surfboard_exporter/pkg/utils"
	log "github.com/sirupsen/logrus"
	// "net/http"
	"net/url"
	"regexp"
	"strconv"
	// "strings"
)

type SB6120 struct {
	address            url.URL
	DownstreamChannels []modem.DownstreamChannel
	ModemInfo          modem.ModemInfo
}

type downstream struct {
	channelID  string
	frequency  float64
	snr        float64
	modulation string
	power      float64
}

type upstream struct {
	channelID        string
	frequency        float64
	rangingServiceID string
	symbolRate       float64
	power            float64
	modulation       string
	rangingStatus    float64
}

type codewords struct {
	channelID     string
	unerrored     float64
	correctable   float64
	uncorrectable float64
}

// DetectModel retruns true if if the modem at the provided address is a SB6120
func DetectModel(address url.URL) bool {
	// 	address.Path = "cmstatushelp.html"
	// 	page := utils.GetPage(address)
	// 	model := page.Find("span#thisModelNumberIs").Text()
	// 	return model == "SB6120"

	return false
}

func New(address url.URL) *SB6120 {
	return &SB6120{
		address: address,
		ModemInfo: modem.ModemInfo{
			Vendor:        "Motorola",
			DocsisVersion: "3.0",
			Model:         "SB6120",
		},
	}
}

// Refresh parses the modem's web interface and updates the class attributes
func (sb *SB6120) Refresh() error {

	err := sb.parseSignalPage()
	if err != nil {
		return err
	}
	// sb.parseHelpPage()
	// sb.parseAddressesPage()

	return nil
}

func (sb *SB6120) GetModemInfo() modem.ModemInfo {
	return sb.ModemInfo
}

// getPage downloads and returns the url document body
// func getPage(url string) *goquery.Document {

// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != 200 {
// 		log.Fatalf("status code error: %s", res.Status)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return doc
// }

func (sb *SB6120) GetDownstreamChannels() []modem.DownstreamChannel {
	var channels []modem.DownstreamChannel
	return channels
}

func (sb *SB6120) GetUpstreamChannels() []modem.UpstreamChannel {
	var channels []modem.UpstreamChannel
	return channels
}

// parseSignalPage parses the modem's signal page for upstream/downstream info
func (sb *SB6120) parseSignalPage() error {

	// Get the document
	sb.address.Path = "cmSignalData.htm"
	page, err := utils.GetPage(sb.address)
	if err != nil {
		return err
	}

	// page, err := utils.GetPage("http://192.168.100.1/cmSignalData.htm")

	// Get the centers on the page (these contain the tables)
	centers := page.Find("html body").ChildrenFiltered("center")

	parseDownstream(centers.Eq(0).ChildrenFiltered("table"))
	// parseUpstream(centers.Eq(1).ChildrenFiltered("table"))
	// parseCodewords(centers.Eq(2).ChildrenFiltered("table"))

	return nil

}

func parseDownstream(s *goquery.Selection) []modem.DownstreamChannel {

	var ds []modem.DownstreamChannel

	// Downstream data is in the first table
	// There are 4 downstream channels to parse. We'll do them at the same time
	// var ds [4]downstream

	channels := s.Find("tr").Eq(1).ChildrenFiltered("td")
	for i := 0; i < 4; i++ {
		value, err := strconv.ParseInt(getFieldValue(channels.Eq(1+i).Text()), 10, 64)
		if err != nil {
			log.Warn(err)
		}

		ds[i].ChannelID = int(value)
	}

	frequency := s.Find("tr").Eq(2).ChildrenFiltered("td")
	for i := 0; i < 4; i++ {
		value, err := strconv.ParseInt(getFieldValue(frequency.Eq(1+i).Text()), 10, 64)
		if err != nil {
			log.Warn(err)
		}

		ds[i].Frequency = value
	}

	snr := s.Find("tr").Eq(3).ChildrenFiltered("td")
	for i := 0; i < 4; i++ {
		value, err := strconv.ParseFloat(getFieldValue(snr.Eq(1+i).Text()), 64)
		if err != nil {
			log.Warn(err)
		}

		ds[i].SignalToNoiseRatio = value
	}

	mod := s.Find("tr").Eq(4).ChildrenFiltered("td")
	for i := 0; i < 4; i++ {
		value := getFieldValue(mod.Eq(1 + i).Text())
		ds[i].Modulation = value
	}

	power := s.Find("tr").Eq(5).ChildrenFiltered("td")
	for i := 0; i < 4; i++ {
		value, err := strconv.ParseFloat(getFieldValue(power.Eq(1+i).Text()), 64)
		if err != nil {
			log.Warn(err)
		}

		ds[i].Power = value
	}

	return ds

}

// func parseUpstream(s *goquery.Selection) {

// 	// Downstream data is in the first table
// 	// There are 4 downstream channels to parse. We'll do them at the same time

// 	var us [4]upstream

// 	channels := s.Find("tr").Eq(1).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value := getFieldValue(channels.Eq(1 + i).Text())
// 		us[i].channelID = value
// 	}

// 	frequency := s.Find("tr").Eq(2).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(frequency.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		us[i].frequency = value
// 	}

// 	rsid := s.Find("tr").Eq(3).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value := getFieldValue(rsid.Eq(1 + i).Text())
// 		us[i].rangingServiceID = value
// 	}

// 	sr := s.Find("tr").Eq(4).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(sr.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		us[i].symbolRate = value
// 	}

// 	power := s.Find("tr").Eq(5).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(power.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		us[i].power = value
// 	}

// 	mod := s.Find("tr").Eq(6).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value := mod.Eq(1 + i).Text()
// 		value = strings.ReplaceAll(value, "<br/>", "")
// 		value = strings.ReplaceAll(value, "\n", ",")
// 		value = strings.Trim(value, ",\u00a0")
// 		us[i].modulation = value
// 	}

// 	status := s.Find("tr").Eq(7).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value := status.Eq(1 + i).Text()
// 		value = strings.Trim(value, "\u00a0")
// 		if value == "Success" {
// 			us[i].rangingStatus = 3
// 		} else if value == "Abort" {
// 			us[i].rangingStatus = 2
// 		} else if value == "Continue" {
// 			us[i].rangingStatus = 1
// 		} else {
// 			// Unknown
// 			us[i].rangingStatus = 0
// 		}
// 	}

// 	for _, u := range us {
// 		modemUpstreamChannelInfo.With(prometheus.Labels{"channel_id": u.channelID, "ranging_service_id": u.rangingServiceID, "modulation_methods": u.modulation}).Set(1)
// 		modemUpstreamFrequency.With(prometheus.Labels{"channel_id": u.channelID}).Set(u.frequency)
// 		modemUpstreamSymbolRate.With(prometheus.Labels{"channel_id": u.channelID}).Set(u.symbolRate)
// 		modemUpstreamPower.With(prometheus.Labels{"channel_id": u.channelID}).Set(u.power)
// 		modemUpstreamRangingStatus.With(prometheus.Labels{"channel_id": u.channelID}).Set(u.rangingStatus)
// 	}

// }

// func parseCodewords(s *goquery.Selection) {

// 	var cw [4]codewords

// 	channels := s.Find("tr").Eq(1).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value := getFieldValue(channels.Eq(1 + i).Text())
// 		cw[i].channelID = value
// 	}

// 	unerrored := s.Find("tr").Eq(2).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(unerrored.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		cw[i].unerrored = value
// 	}

// 	correctable := s.Find("tr").Eq(3).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(correctable.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		cw[i].correctable = value
// 	}

// 	uncorrectable := s.Find("tr").Eq(4).ChildrenFiltered("td")
// 	for i := 0; i < 4; i++ {
// 		value, err := strconv.ParseFloat(getFieldValue(uncorrectable.Eq(1+i).Text()), 64)
// 		if err != nil {
// 			log.Warn(err)
// 		}
// 		cw[i].uncorrectable = value
// 	}

// 	for _, c := range cw {
// 		modemCodewordsUnerrored.With(prometheus.Labels{"channel_id": c.channelID}).Set(c.unerrored)
// 		modemCodewordsCorrectable.With(prometheus.Labels{"channel_id": c.channelID}).Set(c.correctable)
// 		modemCodewordsUncorrectable.With(prometheus.Labels{"channel_id": c.channelID}).Set(c.uncorrectable)
// 	}

// }

// func parseAddressesPage() {

// 	labels := prometheus.Labels{}

// 	// Get the document
// 	page := getPage("http://192.168.100.1/cmAddressData.htm")

// 	// Get the centers on the page (these contain the tables)
// 	centers := page.Find("html body").ChildrenFiltered("center")

// 	// Parse the first table
// 	tablerows := centers.Eq(0).ChildrenFiltered("table").Find("tr")
// 	labels["serial_number"] = tablerows.Eq(1).ChildrenFiltered("td").Eq(1).Text()
// 	labels["hfc_mac"] = tablerows.Eq(2).ChildrenFiltered("td").Eq(1).Text()
// 	labels["ethernet_ip"] = tablerows.Eq(3).ChildrenFiltered("td").Eq(1).Text()
// 	labels["ethernet_mac"] = tablerows.Eq(4).ChildrenFiltered("td").Eq(1).Text()

// 	// Parse the second table
// 	tablerows = centers.Eq(1).ChildrenFiltered("table").Find("tr")
// 	labels["cpe_mac"] = tablerows.Eq(1).ChildrenFiltered("td").Eq(1).Text()
// 	labels["cpe_mac_status"] = tablerows.Eq(1).ChildrenFiltered("td").Eq(2).Text()

// 	modemAddresses.With(labels).Set(1)
// }

// // parseHelpPage parses the modem's help page for hardward/firmware info
// func parseHelpPage() {

// 	// var labels prometheus.Labels
// 	labels := make(map[string]string)

// 	// Get the document
// 	page := getPage("http://192.168.100.1/cmHelpData.htm")

// 	// Get the tables on the page
// 	tables := page.Find("html body table")

// 	// The data we want is in the first table (before JS loads header), first row, first column
// 	html, _ := tables.Eq(0).Find("td").Eq(0).Html()

// 	// The element is a line-break separated string
// 	parts := strings.Split(html, "<br/>")

// 	// This will parse the fields with the info and get the text to the right
// 	// of the first colon (:), trimmed of whitespace
// 	getInfoValue := func(text string) string {
// 		return strings.TrimSpace(text[strings.Index(text, ":")+1 : len(text)])
// 	}

// 	// Create our label fields
// 	labels["model_name"] = getInfoValue(parts[0])
// 	labels["vendor_name"] = getInfoValue(parts[1])
// 	labels["firmware_name"] = getInfoValue(parts[2])
// 	labels["boot_version"] = getInfoValue(parts[3])
// 	labels["hardware_version"] = getInfoValue(parts[4])
// 	labels["serial_number"] = getInfoValue(parts[5])
// 	labels["firmware_build_time"] = getInfoValue(parts[6])

// 	// Set the metric
// 	modemInfo.With(labels).Set(1)
// }

func getFieldValue(text string) string {

	reg, err := regexp.Compile("[0-9A-Za-z.]+")
	if err != nil {
		log.Warn(err)
		return ""
	}
	text = reg.FindString(text)

	return text
}
