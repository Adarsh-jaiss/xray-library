package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// DB_PASSWORD is the environment variable that holds the database password.
var DB_PASSWORD = "DB_PASSWORD"

const (
	// POSTGRES_SCHEMA_QUERY is the SQL query used to describe a table schema in PostgreSQL.
	POSTGRES_SCHEMA_QUERY = "SELECT column_name, data_type, character_maximum_length FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = %s;"

	// POSTGRES_TABLE_LIST_QUERY is the SQL query used to list all tables in a schema in PostgreSQL.
	POSTGRES_TABLE_LIST_QUERY = "SELECT table_name FROM information_schema.tables WHERE table_schema= %s AND table_type='BASE TABLE';"
)

// Postgres is a PostgreSQL implementation of the ISQL interface.
type Postgres struct {
	Client *sql.DB
}

// NewPostgres creates a new PostgreSQL client with the given sql.DB.
func NewPostgres(dbClient *sql.DB) (types.ISQL, error) {
	return &Postgres{
		Client: dbClient,
	}, nil

}

// NewPostgresWithConfig creates a new PostgreSQL client with the given configuration.
// It returns an error if the DB_PASSWORD environment variable is not set.
func NewPostgresWithConfig(dbConfig *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 {
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}
	DB_PASSWORD = os.Getenv(DB_PASSWORD)

	dbtype := types.Postgres
	db, err := sql.Open(dbtype.String(), fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, DB_PASSWORD, dbConfig.DatabaseName, dbConfig.SSL))
	if err != nil {
		return nil, fmt.Errorf("database connecetion failed : %v", err)
	}
	return &Postgres{
		Client: db,
	}, nil
}

// Schema returns the schema of a table in the database.
// It returns an error if the SQL query fails.
func (p *Postgres) Schema(table string) (types.Table, error) {

	var response types.Table

	// execute the sql statement
	rows, err := p.Client.Query(fmt.Sprintf(POSTGRES_SCHEMA_QUERY, table))
	if err != nil {
		return response, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer rows.Close()

	// scanning the result into and append it into a varibale
	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(&column.Name, &column.Type, &column.IsNullable, &column.Key, &column.DefaultValue, &column.Extra, &column.IsPrimary, &column.IsIndex); err != nil {
			return response, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""      // default description
		column.Metatags = []string{} // default metatags as an empty slice
		column.Metatags = append(column.Metatags, column.Name)
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
	if err := rows.Err(); err != nil {
		return response, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil

}

// Execute executes a SQL query and returns the result as a JSON byte slice.
// It returns an error if the SQL query fails.
func (p *Postgres) Execute(query string) ([]byte, error) {
	// execute the sql statement
	rows, err := p.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	// getting the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
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

// Tables returns a list of all tables in the given database.
// It returns an error if the SQL query fails.
func (p *Postgres) Tables(databaseName string) ([]string, error) {
	rows, err := p.Client.Query(fmt.Sprintf(POSTGRES_TABLE_LIST_QUERY, databaseName))
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning database")
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error interating over rows: %v", err)
	}

	return tables, nil

}
