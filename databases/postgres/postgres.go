package postgres

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// DB_PASSWORD is the environment variable that holds the database password.
var DB_PASSWORD = "DB_PASSWORD"

const (
	// POSTGRES_SCHEMA_QUERY is the SQL query used to describe a table schema in PostgreSQL.
	// POSTGRES_SCHEMA_QUERY = "SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1;"
	POSTGRES_SCHEMA_QUERY = `
	SELECT 
    	c.column_name AS name,
    	c.data_type AS type,
    	c.is_nullable AS is_nullable,
    	c.column_default AS default_value,
    	c.character_maximum_length AS character_maximum_length,
    	c.ordinal_position AS ordinal_position,
    	CASE WHEN c.column_default IS NOT NULL THEN true ELSE false END AS visibility,
    	CASE WHEN kcu.column_name IS NOT NULL THEN true ELSE false END AS is_primary,
    	CASE WHEN c.is_updatable = 'YES' THEN true ELSE false END AS is_updatable
	FROM 
    	information_schema.columns c
	LEFT JOIN 
    	information_schema.key_column_usage kcu 
	ON 
    	c.table_name = kcu.table_name AND c.column_name = kcu.column_name
	WHERE 
    	c.table_name = $1;
	`

	// POSTGRES_TABLE_LIST_QUERY is the SQL query used to list all tables in a schema in PostgreSQL.
	POSTGRES_TABLE_LIST_QUERY = "SELECT table_name FROM information_schema.tables WHERE table_schema= 'public' AND table_type='BASE TABLE' AND table_catalog = $1;"
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
	db, err := sql.Open(dbtype.String(), fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, DB_PASSWORD, dbConfig.Database, dbConfig.SSL))
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
	rows, err := p.Client.Query(POSTGRES_SCHEMA_QUERY, table)
	if err != nil {
		return response, fmt.Errorf("error executing sql statement: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// scanning the result into and append it into a variable
	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&column.IsNullable,
			&column.DefaultValue,
			&column.CharacterMaximumLength,
			&column.OrdinalPosition,
			&column.Visibility,
			&column.IsPrimary,
			&column.IsUpdatable,
		); err != nil {
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

	tbl, err := types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil

	fmt.Println(TableToString(tbl))
	return tbl, nil
}

// Execute executes a SQL query and returns the result as a JSON byte slice.
// It returns an error if the SQL query fails.
func (p *Postgres) Execute(query string) ([]byte, error) {
	// execute the sql statement
	query = PostgresMetaCommands(query)
	rows, err := p.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// getting the column names
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
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Decode base64 data
		for _, val := range values {
			strVal, ok := val.(*string)
			if ok && strVal != nil && isBase64(*strVal) {
				decoded, err := base64.StdEncoding.DecodeString(*strVal)
				if err != nil {
					return nil, fmt.Errorf("error decoding base64 data: %v", err)
				}
				*strVal = string(decoded)
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

// Tables returns a list of all tables in the given database.
// It returns an error if the SQL query fails.
func (p *Postgres) Tables(databaseName string) ([]string, error) {
	rows, err := p.Client.Query(POSTGRES_TABLE_LIST_QUERY, databaseName)
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
			return nil, fmt.Errorf("error scanning database")
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error interating over rows: %v", err)
	}

	return tables, nil
}

// GenerateCreateTableQuery generates a CREATE TABLE query for the given table.
// It returns the query as a string.
func (p *Postgres) GenerateCreateTableQuery(table types.Table) string {
	query := `CREATE TABLE "` + table.Name + `" (`
	for i, col := range table.Columns {
		colType := strings.ToUpper(col.Type)
		query += col.Name + " " + colType
		if col.IsPrimary {
			query += " PRIMARY KEY"
		}
		if col.DefaultValue.Valid {
			query += " DEFAULT " + col.DefaultValue.String
		}
		if col.IsUnique.String == "YES" {
			query += " UNIQUE"
		}
		if col.IsNullable == "NO" {
			query += " NOT NULL"
		}
		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ");"
	return query
}

// TableToString returns a string representation of a table.
// It is used for debugging purposes.
func TableToString(t types.Table) string {
	var cols []string
	for _, col := range t.Columns {
		cols = append(cols, fmt.Sprintf(
			"Name: %s, Type: %s, IsNullable: %v, DefaultValue: %v, CharacterMaximumLength: %v, OrdinalPosition: %v, Visibility: %v, IsPrimary: %v, IsUpdatable: %v",
			col.Name,
			col.Type,
			col.IsNullable,
			col.DefaultValue,
			col.CharacterMaximumLength,
			col.OrdinalPosition,
			col.Visibility,
			col.IsPrimary,
			col.IsUpdatable,
		))
	}
	return fmt.Sprintf(
		"Table: %s, Columns: [%s], ColumnCount: %d",
		t.Name,
		strings.Join(cols, "; "),
		t.ColumnCount,
	)
}

func PostgresMetaCommands(query string) string {
	switch query {
	case "\\l":
		return "SELECT datname FROM pg_database WHERE datistemplate = false;"
	case "\\dt":
		return "SELECT * FROM pg_catalog.pg_tables;"
	case "\\d":
		return "SELECT * FROM pg_catalog.pg_tables;"
	case "\\c":
		return "switch_database"
	case "\\q":
		return "exit"
	case "\\?":
		return "help"
	case "\\h":
		return "help"
	case "\\du":
		return "SELECT * FROM pg_catalog.pg_roles;"
	case "\\conninfo":
		return "SELECT * FROM pg_stat_activity WHERE pid = pg_backend_pid();"
	default:
		// Handle meta commands with parameters
		if strings.HasPrefix(query, "\\c ") {
			dbName := strings.TrimPrefix(query, "\\c ")
			return fmt.Sprintf("switch_database %s", dbName)
		} else if strings.HasPrefix(query, "\\d ") {
			tableName := strings.TrimPrefix(query, "\\d ")
			return fmt.Sprintf("SELECT * FROM %s;", tableName)
		} else if strings.HasPrefix(query, "\\dn ") {
			schemaName := strings.TrimPrefix(query, "\\dn ")
			return fmt.Sprintf("SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname = '%s';", schemaName)
		} else if strings.HasPrefix(query, "\\dp ") {
			tableName := strings.TrimPrefix(query, "\\dp ")
			return fmt.Sprintf("SELECT * FROM pg_catalog.pg_statio_all_tables WHERE relname = '%s';", tableName)
		}
	}
	// If the query doesn't match any known meta commands, return it unchanged
	return query
}
