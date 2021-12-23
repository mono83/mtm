# mtm
MySQL Telemetry Monitoring is a simple Go application, that sends data about MySQL table into InfluxDB/Prometheus.

### Data sent:
- `mysql.table.size.total` - total table size (data + index)
- `mysql.table.size.data` - table data size
- `mysql.table.size.index` - table index size
- `mysql.table.rows` - row count 
- `mysql.table.partitions` - partitions, if present

All metrics are shipped with tags `database`,`table` and `engine`

