package bigquery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	_ "gorm.io/driver/bigquery/driver"
)

var GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"

const (
	BigQuery_SCHEMA_QUERY = "SELECT * FROM %s LIMIT 1"
	BigQuery_TABLES_QUERY = "SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%s'"
)

// The BigQuery struct is responsible for holding the BigQuery client and configuration.
type BigQuery struct {
	Client *sql.DB
	Config *config.Config
}

// NewBigQuery creates a new instance of BigQuery with the provided client.
// It returns an instance of types.ISQL and an error.
func NewBigQuery(client *sql.DB) (types.ISQL, error) {
	return &BigQuery{
		Client: client,
		Config: &config.Config{},
	}, nil
}

// NewBigQueryWithConfig creates a new instance of BigQuery with the provided configuration.
// It returns an instance of types.ISQL and an error.
func NewBigQueryWithConfig(cfg *config.Config) (types.ISQL, error) {
	if os.Getenv(GOOGLE_APPLICATION_CREDENTIALS) == "" || len(os.Getenv(GOOGLE_APPLICATION_CREDENTIALS)) == 0 {
		return nil, fmt.Errorf("please set %s env variable for the database", GOOGLE_APPLICATION_CREDENTIALS)
	}
	GOOGLE_APPLICATION_CREDENTIALS = os.Getenv(GOOGLE_APPLICATION_CREDENTIALS)

	dbType := types.BigQuery
	connectionString := fmt.Sprintf("bigquery://%s/%s", cfg.ProjectID, cfg.DatabaseName)
	db, err := sql.Open(dbType.String(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}

	return &BigQuery{
		Client: db,
		Config: cfg,
	}, nil
}

// this function extarcts the schema of a table in BigQuery.
// It takes table name as input and returns a Table struct and an error.
func (b *BigQuery) Schema(table string) (types.Table, error) {

	// execute the sql statement
	rows, err := b.Client.Query(fmt.Sprintf(BigQuery_SCHEMA_QUERY, table))
	if err != nil {
		return types.Table{}, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// scanning the result into a variable and append it into a the slice
	var columns []types.Column
	columnNames, err := rows.Columns()
	if err != nil {
		return types.Table{}, fmt.Errorf("error getting column names: %v", err)
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return types.Table{}, fmt.Errorf("error getting column types: %v", err)
	}
	for i, name := range columnNames {
		var column types.Column
		column.Name = name
		column.Type = columnTypes[i].DatabaseTypeName()
		var ISNullable, _ = columnTypes[i].Nullable()
		column.IsNullable = fmt.Sprintf("%v", ISNullable)
		var length int64
		length, _ = columnTypes[i].Length()
		column.CharacterMaximumLength = sql.NullInt64{
			Int64: length,
			Valid: length != 0,
		}
		columns = append(columns, column)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		Dataset:     b.Config.DatabaseName,
		ColumnCount: int64(len(columns)),
	}, nil

}

// Execute executes a query on BigQuery.
// It takes a query string as input and returns the result as a byte slice and an error.
func (b *BigQuery) Execute(query string) ([]byte, error) {
	rows, err := b.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		// create a slice of values and pointers
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			//  create a slice of pointers to the values
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		results = append(results, values)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Convert the result to JSON
	queryResult := types.QueryResult{
		Columns: columns,
		Rows:    results,
	}

	jsonData, err := json.Marshal(queryResult)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %v", err)
	}

	return jsonData, nil

}

// Tables returns a list of tables in a dataset.
// It takes a dataset name as input and returns a slice of strings and an error.
func (b *BigQuery) Tables(dataset string) ([]string, error) {
	// res, err := b.Client.Query("SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '" + Dataset + "'")

	rows, err := b.Client.Query(fmt.Sprintf(BigQuery_TABLES_QUERY, dataset))
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var tables []string

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning dataset")
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error interating over rows: %v", err)
	}

	return tables, nil
}

// GenerateCreateTableQuery generates a CREATE TABLE query for BigQuery.
func (b *BigQuery) GenerateCreateTableQuery(table types.Table) string {
	query := "CREATE TABLE " + table.Dataset + "." + table.Name + " ("
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + convertTypeToBigQuery(colType)

		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ");"
	return query
}

// convertTypeToBigQuery converts a Data type to a BigQuery SQL Data type.
func convertTypeToBigQuery(dataType string) string {
	// Map column types to BigQuery equivalents
	switch dataType {
	case "INT":
		return "INT64"
	case "VARCHAR(255)", "TEXT":
		return "STRING"
	case "INTEGER":
		return "INT64"
	case "FLOAT":
		return "FLOAT64"
	case "BOOLEAN":
		return "BOOL"
	case "DATE":
		return "DATE"
	case "DATETIME", "TIMESTAMP":
		return "TIMESTAMP"
	// Add more type conversions as needed
	default:
		return dataType
	}
}
