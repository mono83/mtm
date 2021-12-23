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
)

// MainCmd is main cobra command
var MainCmd = &cobra.Command{
	Use:   "mtm",
	Short: "MySQL Telemetry Monitoring",
}

func genericInject(c *cobra.Command, influx, prometheus bool) {
	c.Flags().StringVar(&flagMySQLDSN, "mysql", "root:root@tcp(127.0.0.1:3306)/", "MySQL DSN")
	if influx {
		c.Flags().StringVar(&flagInfluxDBAddr, "influxdb", "", "InfluxDB UDP listening port")
		c.Flags().IntVar(&flagInfluxPacketSizeLimit, "packet", 1024, "Max UDP packet size limit")
	}
	if prometheus {
		c.Flags().StringVar(&flagPrometheusBind, "prometheus", "", "Prometheus listening port")
	}
}

func genericReader() func() (values.MetricCollection, error) {
	return func() (values.MetricCollection, error) {
		return analyze.MySQLDSN(flagMySQLDSN)
	}
}

func init() {
	MainCmd.AddCommand(
		testCmd,
		onceCmd,
		serveCmd,
	)
}
