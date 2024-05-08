package snowflake

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	sf "github.com/snowflakedb/gosnowflake"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

type Snowflake struct {
	Client *sql.DB
	Config *config.Config
}

var DB_PASSWORD string = "DB_PASSWORD"

const (
	SNOWFLAKE_TABLES_LIST_QUERY = "SHOW TERSE TABLES"
	SNOWFLAKE_SCHEMA_QUERY = "SELECT column_name::TEXT, data_type::TEXT FROM information_schema.columns WHERE table_name::TEXT = ?;"
)

// The NewSnowflake function is responsible for creating a new Snowflake object with an initialized database client and configuration.
func NewSnowflake(dbClient *sql.DB) (types.ISQL, error) {
	return &Snowflake{
		Client: dbClient,
		Config: &config.Config{},
	}, nil
}

// The NewSnowflakeWithConfig function is responsible for creating a new Snowflake object with an initialized database client and configuration.
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
		Database:     config.DatabaseName,
		Warehouse:    config.Warehouse,
		Schema:       config.Schema,
		InsecureMode: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating snowflake DSN: %v", err)
	}
	fmt.Println(dsn)
	dbType := types.Snowflake
	db, err := sql.Open(dbType.String(), dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to snowflake database: %v", err)
	}

	return &Snowflake{
		Client: db,
		Config: config,
	}, nil

}

// The Schema function returns the schema of a table in Snowflake.
func (s *Snowflake) Schema(table string) (types.Table, error) {
	var res types.Table

	rows, err := s.Client.Query(SNOWFLAKE_SCHEMA_QUERY, table)
	if err != nil {
		return res, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(&column.Name,&column.Type); err != nil {
			return res, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""      // default description
		column.Metatags = []string{} // default metatags as an empty string slice
		column.Metatags = append(column.Metatags, column.Name)
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
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

// Every table in Snowflake lives "inside" a schema. Every schema lives "inside" a database. It's a hierarchical system.
// The Tables function returns a list of tables in a Snowflake database.
func (s *Snowflake) Tables(DatabaseName string) ([]string, error) {
	query := fmt.Sprintf("USE WAREHOUSE %s", s.Config.Warehouse)

	_, err := s.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}

	rows, err := s.Client.Query(SNOWFLAKE_TABLES_LIST_QUERY)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement and querying tables list: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var col1, col2, col3, col4, col5 string
		if err := rows.Scan(&col1, &col2, &col3, &col4, &col5); err != nil {
			return nil, fmt.Errorf("error scanning database: %v", err)
		}
		table := col2
		tables = append(tables, table)
	}

	// checking for errors in iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows:%v", err)
	}

	return tables, nil
}

// The Execute function executes a query on a Snowflake database and returns the result as a JSON byte slice.
func (s *Snowflake) Execute(query string) ([]byte, error) {
	rows, err := s.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

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
