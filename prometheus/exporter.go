package prometheus

import (
	"fmt"
	"github.com/mono83/mtm/values"
	"github.com/mono83/udpwriter/influxdb"
	"io"
	"net/http"
)

// BuildGaugeExporterFunc build HTTP handler func that will
// output metrics gauge data in Prometheus format
func BuildGaugeExporterFunc(provider func() values.MetricCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteGaugeTo(w, provider())
	}
}

// WriteGaugeTo writes data in Prometheus gauge format to given writer
func WriteGaugeTo(w io.Writer, data values.MetricCollection) {
	data.ForEachMetric(func(name string, value int64, tags map[string]string) {
		w.Write([]byte("# TYPE "))
		w.Write(influxdb.Sanitize(name))
		w.Write([]byte(" gauge\n"))

		w.Write(influxdb.Sanitize(name))
		if len(tags) > 0 {
			w.Write([]byte(" {"))
			first := true
			for k, v := range tags {
				if first {
					first = false
				} else {
					w.Write([]byte(","))
				}
				w.Write(influxdb.Sanitize(k))
				w.Write([]byte("=\""))
				w.Write(influxdb.Sanitize(v))
				w.Write([]byte("\""))
			}
			w.Write([]byte("}"))
		}
		w.Write([]byte(" "))
		fmt.Fprintln(w, value)
	})
}
