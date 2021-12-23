package analyze

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/mono83/mtm/values"
	"log"
)

// MySQL performs database server health analysis using connection
// built on given configuration
func MySQL(c mysql.Config) (values.MetricCollection, error) {
	return MySQLDSN(c.FormatDSN())
}

// MySQLDSN performs database server health analysis using connection
// built on given DSN
func MySQLDSN(dsn string) (values.MetricCollection, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return Analyze(db)
}

// Analyze performs database server health analysis using given connection
func Analyze(conn *sql.DB) (values.MetricCollection, error) {
	if conn == nil {
		return nil, errors.New("nil database connection")
	}

	dbs, err := listdb(conn)
	if err != nil {
		return nil, err
	}
	log.Printf("Got %d databases\n", len(dbs))
	var totalTables []table
	for _, d := range dbs {
		log.Printf("Entering database %s\n", d)
		tables, err := listtables(conn, d)
		if err != nil {
			return nil, err
		}
		log.Printf("Got %d tables in %s\n", len(tables), d)
		totalTables = append(totalTables, tables...)
	}

	return analysisResult(totalTables), nil
}

type analysisResult []table

func (a analysisResult) ForEachMetric(f func(name string, value int64, tags map[string]string)) {
	if f == nil {
		return
	}

	for _, t := range a {
		tags := map[string]string{
			"database": t.database,
			"table":    t.name,
			"engine":   t.engine,
		}

		f("mysql.table.size.total", int64(t.data+t.index), tags)
		f("mysql.table.size.data", int64(t.data), tags)
		f("mysql.table.size.index", int64(t.index), tags)
		f("mysql.table.rows", int64(t.rows), tags)

		if len(t.partitions) > 0 {
			f("mysql.table.partitions", int64(len(t.partitions)), tags)
		}
	}
}
