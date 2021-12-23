package cmd

import (
	"github.com/mono83/udpwriter/influxdb"
	"github.com/spf13/cobra"
	"time"
)

var onceCmd = &cobra.Command{
	Use:     "once",
	Aliases: []string{"influx", "influxdb"},
	Short:   "Reads data and flushes it to InfluxDB, then quits",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := genericReader()()
		if err != nil {
			return err
		}

		w, err := influxdb.NewUDPWriterS(flagInfluxDBAddr, flagInfluxPacketSizeLimit)
		if err != nil {
			return err
		}
		err = w.Write(data)
		time.Sleep(time.Second)
		return err
	},
}

func init() {
	genericInject(onceCmd, true, false)
}
