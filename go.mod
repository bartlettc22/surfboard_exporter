module github.com/bartlettc22/surfboard_exporter

go 1.14

replace github.com/bartlettc22/surfboard_exporter/pkg => ./pkg

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
