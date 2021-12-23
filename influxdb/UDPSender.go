package influxdb

import (
	"github.com/mono83/mtm/values"
	"github.com/mono83/udpwriter/influxdb"
)

// SendUDP sends metrics using UDP to InfluxDB server
func SendUDP(addr string, buf int, collection values.MetricCollection) error {
	w, err := influxdb.NewUDPWriterS(addr, buf)
	if err != nil {
		return err
	}
	return w.Write(collection)
}
