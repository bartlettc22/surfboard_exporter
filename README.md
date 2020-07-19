# Surfboard Cable Modem Scraper and Prometheus Exporter
[![Build][Build-Status-Image]][Build-Status-Url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]
This exporter scrapes Surfboard modem web interfaces and presents the metrics in Prometheus format.

This project contains the following:
* Prometheus exporter for Surfboard modems (amd64/arm/Docker)
* Golang library for interfacing with the Surfboard cable modems

## Currently Supported Modems
* Motorola SB6120 Cable Modem
* Arris SB8200 Cable Modem

## Using the Prometheus exporter
To use the Prometheus exporter, simply download the latest binary from the [releases page](https://github.com/bartlettc22/surfboard_exporter/releases/latest) and run it
```
./surfboard_exporter
```

The following arguments are available

|Argument|Default|Description|
|-|-|-|
|`--listen-port`, `-p`|`9040`|Port that metrics server listens on. Metrics available at (`<host>:<listen-port>/metrics`)|
|`--modem-address`, `-a`|`http://192.168.100.1`|URL address of the modem|
|`--modem-model`, `-m`|`auto`|Model of modem [auto, sb6120, sb8200]|

The exporter can also be run in Docker like so
```
docker run -d -p 9040:9040 bartlettc/surfboard_exporter:v0.2.0 -p 9040
```

<!--
## Usage
```
surfboard_exporter --address 192.168.100.1 --port 9040
```

| Parameter | Description |
|---|---|---|---|---|
| --address | Address of the modem's web interface |
| --port -p | Metrics port |
-->

[Build-Status-Url]: https://travis-ci.org/bartlettc22/surfboard_exporter
[Build-Status-Image]: https://travis-ci.org/bartlettc22/surfboard_exporter.svg?branch=master
[reportcard-url]: https://goreportcard.com/report/github.com/bartlettc22/surfboard_exporter
[reportcard-image]: https://goreportcard.com/badge/github.com/bartlettc22/surfboard_exporter
[godoc-url]: https://godoc.org/github.com/bartlettc22/surfboard_exporter/pkg
[godoc-image]: https://godoc.org/github.com/bartlettc22/surfboard_exporter/pkg?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg