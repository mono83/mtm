package values

// MetricCollection defines container for metrics data
type MetricCollection interface {
	ForEachMetric(func(name string, value int64, tags map[string]string))
}
