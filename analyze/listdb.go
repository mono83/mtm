package analyze

import "database/sql"

func listdb(conn *sql.DB) ([]string, error) {
	rows, err := conn.Query("SELECT `SCHEMA_NAME` FROM `information_schema`.`schemata`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			return nil, err
		}
		databases = append(databases, row)
	}

	return databases, nil
}
