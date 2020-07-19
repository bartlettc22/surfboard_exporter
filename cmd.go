package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "surfboard_exporter",
	Short: "Prometheus metrics exporter for Surfboard cable modems",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of surfboard_exporter",
	Long:  `All software has versions. This is surfboard_exporter's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("surfboard_exporter v%s\n", version)
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&listenPort, "listen-port", "p", 9040, "Port that metrics server listens on. Metrics available at (`<host>:<listen-port>/metrics`)")
	rootCmd.PersistentFlags().StringVarP(&modemAddress, "modem-address", "a", "http://192.168.100.1", "URL address of the modem")
	rootCmd.PersistentFlags().StringVarP(&modemModel, "modem-model", "m", "auto", "Model of modem [auto, sb6120, sb8200]")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
