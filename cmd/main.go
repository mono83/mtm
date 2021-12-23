package cmd

import (
	"github.com/mono83/mtm/analyze"
	"github.com/mono83/mtm/values"
	"github.com/spf13/cobra"
)

var (
	flagMySQLDSN              string
	flagInfluxDBAddr          string
	flagInfluxPacketSizeLimit int
	flagPrometheusBind        string
	flagFilterOnly            []string
	flagFilterExclude         []string
)

// MainCmd is main cobra command
var MainCmd = &cobra.Command{
	Use:   "mtm",
	Short: "MySQL Telemetry Monitoring",
}

func genericInject(c *cobra.Command, influx, prometheus bool) {
	c.Flags().StringVar(&flagMySQLDSN, "mysql", "root:root@tcp(127.0.0.1:3306)/", "MySQL DSN")
	c.Flags().StringSliceVar(&flagFilterOnly, "only", nil, "White list databases")
	c.Flags().StringSliceVar(&flagFilterExclude, "exclude", nil, "Black list databases")
	if influx {
		c.Flags().StringVar(&flagInfluxDBAddr, "influxdb", "", "InfluxDB UDP listening port")
		c.Flags().IntVar(&flagInfluxPacketSizeLimit, "packet", 1024, "Max UDP packet size limit")
	}
	if prometheus {
		c.Flags().StringVar(&flagPrometheusBind, "prometheus", "", "Prometheus listening port")
	}
}

func genericReader() func() (values.MetricCollection, error) {
	allow := make(map[string]bool)
	deny := make(map[string]bool)

	for _, f := range flagFilterOnly {
		allow[f] = true
	}
	for _, f := range flagFilterExclude {
		deny[f] = true
	}

	return func() (values.MetricCollection, error) {
		return analyze.MySQLDSN(flagMySQLDSN, func(db string) bool {
			if len(deny) > 0 {
				if _, ok := deny[db]; ok {
					// Matched in blacklist
					return false
				}
			}
			if len(allow) > 0 {
				if _, ok := allow[db]; ok {
					return true
				}
				return false
			}
			return true
		})
	}
}

func init() {
	MainCmd.AddCommand(
		testCmd,
		onceCmd,
		serveCmd,
	)
}
