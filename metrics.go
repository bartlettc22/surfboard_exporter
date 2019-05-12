package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	modemInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "modem_info",
		Help: "Information about the modem hardware/firmware",
	}, []string{"model_name", "vendor_name", "firmware_name", "boot_version", "hardware_version", "serial_number", "firmware_build_time"})
)
