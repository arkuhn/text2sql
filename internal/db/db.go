package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func ensureSSLMode(connectionString string) string {
	if !strings.Contains(connectionString, "sslmode=") {
		if strings.Contains(connectionString, "?") {
			connectionString += "&sslmode=disable"
		} else {
			connectionString += "?sslmode=disable"
		}
	}
	return connectionString
}

func GetTableNames(connectionString string) ([]string, error) {
	connectionString = ensureSSLMode(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, fmt.Errorf("error querying table names: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("error scanning table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func GetTableSchema(connectionString, tableName string) ([]string, error) {
	connectionString = ensureSSLMode(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("error querying table schema: %v", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return nil, fmt.Errorf("error scanning column info: %v", err)
		}
		columns = append(columns, fmt.Sprintf("%s (%s)", columnName, dataType))
	}

	return columns, nil
}

func ExecuteQuery(connectionString, query string) ([]map[string]interface{}, error) {
	connectionString = ensureSSLMode(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column names: %v", err)
	}

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}

		result = append(result, row)
	}

	return result, nil
}

