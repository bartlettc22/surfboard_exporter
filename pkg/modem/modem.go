package modem

import (
	"time"
)

// ModemInfo represents the metadata info for the modem
type ModemInfo struct {
	Vendor            string
	Model             string
	DocsisVersion     string
	HardwareVersion   string
	SoftwareVersion   string
	MACAddress        string
	SerialNumber      string
	FirmwareName      string
	FirmwareBuildTime string
	BootVersion       string
	Uptime            time.Duration
}

// DownstreamChannel represents a single downstream channel signal
type DownstreamChannel struct {
	ChannelID          int
	Locked             bool
	Frequency          int64   // Hz
	SignalToNoiseRatio float64 // dB
	Power              float64 // dBmV
	Modulation         string
	Corrected          int64
	Uncorrectables     int64
}

// DownstreamChannel represents a single downstream channel signal
type UpstreamChannel struct {
	ChannelID         int
	Channel           int
	Locked            bool
	USChannelType     string
	Frequency         int64   // Hz
	Width             int64   // Hz
	Power             float64 // dBmV
	RangingServiceID  string
	ModulationMethods string
}
