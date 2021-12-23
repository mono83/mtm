package analyze

import (
	"database/sql"
	"strings"
)

type table struct {
	database  string
	name      string
	engine    string
	format    string
	collation string

	rows  uint64
	data  uint64
	free  uint64
	index uint64

	partitions []partition
}

type partition struct {
	database string
	table    string
	name     string
	method   string

	rows  uint64
	data  uint64
	free  uint64
	index uint64
}

func listtables(conn *sql.DB, database string) ([]table, error) {
	if "information_schema" == strings.ToLower(database) {
		// Self analysis is forbidden
		return nil, nil
	}
	if "sys" == strings.ToLower(database) {
		// Self analysis is forbidden
		return nil, nil
	}

	// Reading partitions
	partitions, err := listpartitions(conn, database)
	if err != nil {
		return nil, err
	}

	// Reading tables
	rows, err := conn.Query(
		"SELECT `TABLE_SCHEMA`,`TABLE_NAME`,`ENGINE`,`ROW_FORMAT`,`TABLE_COLLATION`,`TABLE_ROWS`,`DATA_LENGTH`,`INDEX_LENGTH`,`DATA_FREE`"+
			" FROM `information_schema`.`tables` WHERE `TABLE_SCHEMA` = ?",
		database,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []table
	for rows.Next() {
		var row table
		err := rows.Scan(
			&row.database,
			&row.name,
			&row.engine,
			&row.format,
			&row.collation,
			&row.rows,
			&row.data,
			&row.index,
			&row.free,
		)
		if err != nil {
			return nil, err
		}
		key := [2]string{row.database, row.name}
		row.partitions = partitions[key]
		tables = append(tables, row)
	}
	return tables, nil
}

func listpartitions(conn *sql.DB, database string) (map[[2]string][]partition, error) {
	rows, err := conn.Query(
		"SELECT `TABLE_SCHEMA`,`TABLE_NAME`,`PARTITION_NAME`,`PARTITION_METHOD`,`TABLE_ROWS`,`DATA_LENGTH`,`INDEX_LENGTH`,`DATA_FREE`"+
			" FROM `information_schema`.`partitions` WHERE `TABLE_SCHEMA` = ? AND `PARTITION_NAME` IS NOT NULL",
		database,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := map[[2]string][]partition{}
	for rows.Next() {
		var row partition
		err := rows.Scan(
			&row.database,
			&row.table,
			&row.name,
			&row.method,
			&row.rows,
			&row.data,
			&row.index,
			&row.free,
		)
		if err != nil {
			return nil, err
		}
		key := [2]string{row.database, row.table}
		data[key] = append(data[key], row)
	}

	return data, nil
}
