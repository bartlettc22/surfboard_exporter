package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (

	// Info metrics

	modemInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_info",
		Help: "Information about the modem hardware/firmware",
	}, []string{"model_name", "vendor_name", "firmware_name", "boot_version", "hardware_version", "serial_number", "firmware_build_time"})

	// Downstream metrics

	modemDownstreamFrequency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_freq",
		Help: "Modem's downstream signal frequency (Hz)",
	}, []string{"channel_id"})

	modemDownstreamSNR = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_snr",
		Help: "Modem's downstream signal-to-noise ratio (dB)",
	}, []string{"channel_id"})

	modemDownstreamModulation = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_modulation",
		Help: "Modem's downstream modulation method",
	}, []string{"channel_id", "modulation_method"})

	modemDownstreamPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_power",
		Help: "Modem's downstream power level (dBmV)",
	}, []string{"channel_id"})

	// Upstream metrics

	modemUpstreamChannelInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_channel_info",
		Help: "Modem's upstream channel info",
	}, []string{"channel_id","ranging_service_id","modulation_methods"})

	modemUpstreamFrequency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_freq",
		Help: "Modem's upstream signal frequency (Hz)",
	}, []string{"channel_id"})

	modemUpstreamSymbolRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_symbol_rate",
		Help: "Modem's upstream symbol rate (Msym/sec)",
	}, []string{"channel_id"})

	modemUpstreamPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_power",
		Help: "Modem's upstream power level (dBmV)",
	}, []string{"channel_id"})

	modemUpstreamRangingStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_ranging_status",
		Help: "Modem's upstream ranging status",
	}, []string{"channel_id"})
)
