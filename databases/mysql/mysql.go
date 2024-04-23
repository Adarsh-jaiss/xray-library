package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	// "github.com/joho/godotenv"
)

// DB_PASSWORD is the environment variable that stores the database password.
var DB_PASSWORD = "DB_PASSWORD"

const (
	SCHEMA_QUERY            = "DESCRIBE %s"                                                              // SCHEMA_QUERY is the SQL query used to describe a table schema.
	MYSQL_TABLES_LIST_QUERY = "SELECT table_name FROM information_schema.tables WHERE table_schema = ?" // MYSQL_TABLES_LIST_QUERY is the SQL query used to list all tables in a schema.
)

// MySQL is a MySQL implementation of the ISQL interface.
type MySQL struct {
	Client *sql.DB // Client is the MySQL database client.
}

func NewMySQL(dbClient *sql.DB) (types.ISQL, error) {
	return &MySQL{
		Client: dbClient,
	}, nil

}

// NewMySQLWithConfig creates a new MySQL client with the given configuration.
// It returns an error if the DB_PASSWORD environment variable is not set.
func NewMySQLWithConfig(dbConfig *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 { // added mysql to be more verbose about the db type
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}

	DB_PASSWORD = os.Getenv(DB_PASSWORD)

	dsn := dbURLMySQL(dbConfig)

	dbtype := types.MySQL
	db, err := sql.Open(dbtype.String(), dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %v", err)
	}

	return &MySQL{
		Client: db,
	}, nil

}

// This method will accept a table name as input and return the table schema (structure).
func (m *MySQL) Schema(table string) (types.Table, error) {
	var response types.Table

	// execute the sql statement
	rows, err := m.Client.Query(fmt.Sprintf(SCHEMA_QUERY, table))
	if err != nil {
		return response, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer rows.Close()

	// scanning the result into and append it into a varibale
	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(&column.Name, &column.Type, &column.IsNullable, &column.Key, &column.DefaultValue, &column.Extra); err != nil {
			return response, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""      // default description
		column.Metatags = []string{} // default metatags as an empty string slice
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

// Execute a database query and return the result in JSON format
func (m *MySQL) Execute(query string) ([]byte, error) {

	// execute the sql statement
	rows, err := m.Client.Query(query)
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

// Retrieve the names of tables in the specified database.
func (m *MySQL) Tables(databaseName string) ([]string, error) {

	// execute the sql statement
	rows, err := m.Client.Query(MYSQL_TABLES_LIST_QUERY, databaseName)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	//scan and append the result
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning database: %v", err)
		}
		tables = append(tables, table)
	}

	// checking for errors in iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows:%v", err)
	}

	return tables, nil
}

func dbURLMySQL(dbConfig *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		DB_PASSWORD,
		dbConfig.Host,
		dbConfig.DatabaseName,
		dbConfig.SSL,
	)
}
