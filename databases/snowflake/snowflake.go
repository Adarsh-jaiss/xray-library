package snowflake

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	sf "github.com/snowflakedb/gosnowflake"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// Snowflake is a Snowflake implementation of the ISQL interface.
type Snowflake struct {
	Client *sql.DB        // Client is the database client for Snowflake.
	Config *config.Config // Config is the configuration for Snowflake.
}

// DB_PASSWORD is the environment variable name for the database password.
var DB_PASSWORD string = "DB_PASSWORD"

const (
	// SNOWFLAKE_TABLES_LIST_QUERY is the query to list tables in Snowflake.
	SNOWFLAKE_TABLES_LIST_QUERY = "SELECT table_name FROM %s.information_schema.tables WHERE table_schema = '%s';"
	// SNOWFLAKE_SCHEMA_QUERY is the query to retrieve schema information for a table in Snowflake.
	SNOWFLAKE_SCHEMA_QUERY = "SELECT column_name::TEXT, data_type::TEXT FROM information_schema.columns WHERE table_name::TEXT = ?;"
)

// NewSnowflake creates a new Snowflake object with an initialized database client and configuration.
func NewSnowflake(dbClient *sql.DB) (types.ISQL, error) {
	return &Snowflake{
		Client: dbClient,
		Config: &config.Config{},
	}, nil
}

// NewSnowflakeWithConfig creates a new Snowflake object with an initialized database client and configuration.
func NewSnowflakeWithConfig(config *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 {
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}
	DB_PASSWORD = os.Getenv(DB_PASSWORD)

	port, _ := strconv.Atoi(config.Port)

	dsn, err := sf.DSN(&sf.Config{
		Account:      config.Account,
		User:         config.Username,
		Password:     DB_PASSWORD,
		Port:         port,
		Database:     config.Database,
		Warehouse:    config.Warehouse,
		Schema:       config.Schema,
		InsecureMode: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating snowflake DSN: %v", err)
	}
	dbType := types.Snowflake
	db, err := sql.Open(dbType.String(), dsn) // open a connection to the snowflake database
	if err != nil {
		return nil, fmt.Errorf("error opening connection to snowflake database: %v", err)
	}

	return &Snowflake{
		Client: db,
		Config: config,
	}, nil

}

// Schema returns the schema of a table in Snowflake.
// It takes the table name as an argument and returns a Table struct and an error if any.
func (s *Snowflake) Schema(table string) (types.Table, error) {
	var res types.Table

	rows, err := s.Client.Query(SNOWFLAKE_SCHEMA_QUERY, table)
	if err != nil {
		return res, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(&column.Name, &column.Type); err != nil {
			return res, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""      // default description
		column.Metatags = []string{} // default metatags as an empty string slice
		column.Metatags = append(column.Metatags, column.Name)
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for errors from iterating over the rows
	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil
}

// Tables returns a list of tables in a Snowflake database.
// It takes the database name as an argument and returns a slice of table names.
func (s *Snowflake) Tables(databaseName string) ([]string, error) {

	query := fmt.Sprintf(SNOWFLAKE_TABLES_LIST_QUERY, databaseName, s.Config.Schema)
	rows, err := s.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement and querying tables list: %v", err)
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

// Execute executes a query on a Snowflake database and returns the result as a JSON byte slice.
func (s *Snowflake) Execute(query string) ([]byte, error) {
	rows, err := s.Client.Query(query)
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
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Decode base64 data
		for i, val := range values {
			strVal, ok := val.(string)
			if ok && isBase64(strVal) {
				// Redecode the value to get the decoded result
				decoded, err := base64.StdEncoding.DecodeString(strVal)
				if err != nil {
					return nil, fmt.Errorf("error decoding base64 data: %v", err)
				}
				values[i] = string(decoded)
			}
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

// isBase64 checks if a string is a valid base64 string.
func isBase64(s string) bool {
	if len(s)%4 != 0 {
		return false
	}
	// Try to decode the string
	_, err := base64.StdEncoding.DecodeString(s)
	// If decoding succeeds, err will be nil, and the function will return true
	// If decoding fails, err will not be nil, and the function will return false
	// Also we do not have access to decoded value, so we are not using it
	return err == nil
}

// GenerateCreateTableQuery generates a CREATE TABLE query for Snowflake.
// It takes a Table struct as an argument and returns the query as a string.
func (s *Snowflake) GenerateCreateTableQuery(table types.Table) string {
	query := "CREATE TABLE " + table.Name + " ("
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + colType
		if column.AutoIncrement {
			query += " AUTOINCREMENT"
		}
		if column.IsPrimary {
			query += " PRIMARY KEY"
		}
		if column.DefaultValue.Valid {
			query += " DEFAULT " + column.DefaultValue.String
		}
		if column.IsUnique.String == "YES" {
			query += " UNIQUE"
		}
		if column.IsNullable == "NO" && !column.IsPrimary {
			query += " NOT NULL"
		}
		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ");"
	return query
}
