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
	}, []string{"model_name", "vendor_name", "firmware_name", "boot_version", "hardware_version", "serial_number", "firmware_build_time", "docsis_version"})

	modemAddresses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_addresses",
		Help: "Information about the modem addresses",
	}, []string{"serial_number", "hfc_mac", "ethernet_ip", "ethernet_mac", "cpe_mac", "cpe_mac_status"})

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

	modemDownstreamLocked = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_locked",
		Help: "Modem's downstream lock status (0=not locked, 1=locked)",
	}, []string{"channel_id"})

	modemDownstreamCorrected = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_corrected_total",
		Help: "Modem's total downstream corrected packets",
	}, []string{"channel_id"})

	modemDownstreamUncorrectable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_downstream_uncorrectable_total",
		Help: "Modem's total downstream uncorrectable packets",
	}, []string{"channel_id"})

	// Upstream metrics

	modemUpstreamChannelInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_channel_info",
		Help: "Modem's upstream channel info",
	}, []string{"channel_id", "channel", "ranging_service_id", "modulation_methods", "us_channel_type"})

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

	// SB8200 only
	modemUpstreamLocked = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_locked",
		Help: "Modem's upstream lock status (0=not locked, 1=locked)",
	}, []string{"channel_id"})

	// SB8200 only
	modemUpstreamWidth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_upstream_width",
		Help: "Modem's upstream width frequency (Hz)",
	}, []string{"channel_id"})

	// Codeword metrics

	modemCodewordsUnerrored = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_codewords_unerrored_total",
		Help: "Modem's total unerrored codewords",
	}, []string{"channel_id"})

	modemCodewordsCorrectable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_codewords_correctable_total",
		Help: "Modem's total correctable codewords",
	}, []string{"channel_id"})

	modemCodewordsUncorrectable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_codewords_uncorrectable_total",
		Help: "Modem's total uncorrectable codewords",
	}, []string{"channel_id"})
)
