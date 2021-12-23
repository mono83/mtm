package cmd

import (
	"errors"
	"github.com/mono83/mtm/influxdb"
	"github.com/mono83/mtm/prometheus"
	"github.com/mono83/mtm/values"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

var flagRefresh time.Duration

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"server"},
	Short:   "Runs server. InfluxDB and/or Prometheus configuration must be provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(flagInfluxDBAddr) == 0 && len(flagPrometheusBind) == 0 {
			return errors.New("both influxdb and prometheus setting are missing")
		}

		// Creating reader func
		reader := genericReader()

		// Obtaining initial data
		data, err := reader()
		if err != nil {
			return err
		}

		var onUpdate func()

		// Starting refresh process
		go func() {
			for {
				time.Sleep(flagRefresh)
				log.Println("Doing update")
				data, err = reader()
				if err != nil {
					log.Println(err)
				} else if onUpdate != nil {
					onUpdate()
				}
			}
		}()

		if len(flagInfluxDBAddr) > 0 {
			log.Printf("Staring sender to InfluxDB %s\n", flagInfluxDBAddr)
			onUpdate = func() {
				if err := influxdb.SendUDP(flagInfluxDBAddr, flagInfluxPacketSizeLimit, data); err != nil {
					log.Println(err)
				}
			}
		}

		if len(flagPrometheusBind) > 0 {
			// Staring Prometheus exporter
			log.Printf("Staring Prometheus listener on %s\n", flagPrometheusBind)
			return http.ListenAndServe(
				flagPrometheusBind,
				httpHandler(prometheus.BuildGaugeExporterFunc(func() values.MetricCollection { return data })),
			)
		} else {
			select {}
		}
	},
}

func init() {
	genericInject(serveCmd, true, true)
	serveCmd.Flags().DurationVar(&flagRefresh, "refresh", time.Minute*15, "Refresh interval")
}

type httpHandler http.HandlerFunc

func (h httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h(w, req)
}
