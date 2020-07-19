package main

import (
	"fmt"
	m "github.com/bartlettc22/surfboard_exporter/pkg/modem"
	"github.com/bartlettc22/surfboard_exporter/pkg/sb6120"
	"github.com/bartlettc22/surfboard_exporter/pkg/sb8200"
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
)

// Generic modem interface
type Modem interface {
	Refresh() error
	GetDownstreamChannels() []m.DownstreamChannel
	GetUpstreamChannels() []m.UpstreamChannel
	GetModemInfo() m.ModemInfo
}

// MetricsHandler is our http request handler
type MetricsHandler struct{}

// var metricsHandler MetricsHandler
var version string
var listenPort int
var modemAddress string
var modemModel string
var modem Modem

// var modem *modem

func init() {
	// Version should be set on build, but if not, set a default
	if version == "" {
		version = "0.0.0"
	}
}

func main() {
	rootCmd.Execute()
}

// Starts the scraping and serving of metrics
func start() {

	var metricsHandler MetricsHandler

	// Detects the modem version and returns the appropriate metrics handler
	modem = detectModem()

	// Perform a refresh to provide immediate feedback if modem is unreachable ()
	err := modem.Refresh()
	if err != nil {
		log.Warnf("Error scraping modem web interface: %v", err)
	}

	http.Handle("/metrics", metricsHandler)

	log.Infof("Listening on port %d", listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

func detectModem() Modem {

	var modem Modem

	u, err := url.Parse(modemAddress)
	if err != nil {
		log.Fatal(err)
	}

	if modemModel == "sb8200" || (modemModel == "auto" && sb8200.DetectModel(*u)) {
		log.Info("Scraping for modem model SB8200")
		modem = sb8200.New(*u)
		return modem
	}

	if modemModel == "sb6120" || (modemModel == "auto" && sb6120.DetectModel(*u)) {
		log.Info("Scraping for modem model SB6120")
		modem = sb6120.New(*u)
		return modem
	}

	if modemModel == "auto" {
		log.Fatalf("Could not detect modem model at address: %s", modemAddress)
	} else {
		log.Fatalf("Invalid modem model specified: %s", modemModel)
	}

	return nil

}

// ServeHTTP scrapes the modem data and serves up the Prometheus metrics
func (m MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := modem.Refresh()
	if err != nil {
		log.Warnf("Error scraping modem web interface: %v", err)
	}

	// We reset the metrics because they have labels that may have changed and we want to get rid of the old ones
	modemInfo.Reset()
	modemDownstreamFrequency.Reset()
	modemDownstreamSNR.Reset()
	modemDownstreamPower.Reset()
	modemDownstreamModulation.Reset()
	modemDownstreamLocked.Reset()
	modemUpstreamFrequency.Reset()

	if err == nil {
		mi := modem.GetModemInfo()
		modemInfo.With(prometheus.Labels{
			"model_name":          mi.Model,
			"vendor_name":         mi.Vendor,
			"hardware_version":    mi.HardwareVersion,
			"serial_number":       mi.SerialNumber,
			"firmware_name":       mi.FirmwareName,
			"boot_version":        mi.BootVersion,
			"firmware_build_time": mi.FirmwareBuildTime,
			"docsis_version":      mi.DocsisVersion,
		}).Set(1)

		ds := modem.GetDownstreamChannels()
		for _, d := range ds {
			modemDownstreamFrequency.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(float64(d.Frequency))
			modemDownstreamSNR.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(d.SignalToNoiseRatio)
			modemDownstreamPower.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(d.Power)
			modemDownstreamModulation.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID), "modulation_method": d.Modulation}).Set(1)
			lock := float64(0)
			if d.Locked {
				lock = 1
			}
			modemDownstreamLocked.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(lock)
			modemDownstreamCorrected.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(float64(d.Corrected))
			modemDownstreamUncorrectable.With(prometheus.Labels{"channel_id": strconv.Itoa(d.ChannelID)}).Set(float64(d.Uncorrectables))
		}

		us := modem.GetUpstreamChannels()
		for _, u := range us {
			modemUpstreamChannelInfo.With(prometheus.Labels{
				"channel":            strconv.Itoa(u.Channel),
				"channel_id":         strconv.Itoa(u.ChannelID),
				"ranging_service_id": u.RangingServiceID,
				"modulation_methods": u.ModulationMethods,
				"us_channel_type":    u.USChannelType,
			}).Set(1)
			modemUpstreamFrequency.With(prometheus.Labels{"channel_id": strconv.Itoa(u.ChannelID)}).Set(float64(u.Frequency))
			modemUpstreamPower.With(prometheus.Labels{"channel_id": strconv.Itoa(u.ChannelID)}).Set(u.Power)
			if mi.Model == "SB8200" {
				lock := float64(0)
				if u.Locked {
					lock = 1
				}
				modemUpstreamLocked.With(prometheus.Labels{"channel_id": strconv.Itoa(u.ChannelID)}).Set(lock)
				modemUpstreamWidth.With(prometheus.Labels{"channel_id": strconv.Itoa(u.ChannelID)}).Set(float64(u.Width))
			} else if mi.Model == "SB6120" {

			}
		}
	}

	spew.Dump(1)

	// Let promhttp serve up the metrics page
	promHandler := promhttp.Handler()
	promHandler.ServeHTTP(w, r)
}
