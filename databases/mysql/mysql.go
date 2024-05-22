package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	// "github.com/joho/godotenv"
)

// DB_PASSWORD is the environment variable that stores the database password.
var DB_PASSWORD = "DB_PASSWORD"

const (
	SCHEMA_QUERY            = "DESCRIBE %s"                                                             // SCHEMA_QUERY is the SQL query used to describe a table schema.
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

// Schema retrieves the table schema for the given table name.
// It takes the table name as an argument and returns the table schema as a types.Table object.
func (m *MySQL) Schema(table string) (types.Table, error) {
	var response types.Table

	// execute the sql statement
	rows, err := m.Client.Query(fmt.Sprintf(SCHEMA_QUERY, table))
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

// Execute executes the given SQL query and returns the result as JSON.
// It takes the SQL query as an argument.
func (m *MySQL) Execute(query string) ([]byte, error) {

	// execute the sql statement
	rows, err := m.Client.Query(query)
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
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Convert the values to the appropriate types
		stringRow := make([]interface{}, len(values))
		for i, v := range values {
			switch value := v.(type) {
			case []byte:
				stringRow[i] = string(value)
			case string:
				stringRow[i] = value
			default:
				stringRow[i] = fmt.Sprintf("%v", value)
			}
		}

		// Append the modified row to the results
		results = append(results, stringRow)
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

// Tables retrieves the list of tables in the given database.
// It takes the database name as an argument and returns a list of table names.
func (m *MySQL) Tables(databaseName string) ([]string, error) {

	// execute the sql statement
	rows, err := m.Client.Query(MYSQL_TABLES_LIST_QUERY, databaseName)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// scan and append the result
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

// GenerateCreateTableQuery generates a SQL query to create a table with the same structure as the input table.
func (m *MySQL) GenerateCreateTableQuery(table types.Table) string {
	query := "CREATE TABLE " + table.Name + " ("
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + convertTypeToMysql(colType)
		if column.AutoIncrement {
			query += " AUTO_INCREMENT"
		}
		if column.IsPrimary {
			query += " PRIMARY KEY"
		}
		if column.DefaultValue.Valid {
			query += " DEFAULT " + column.DefaultValue.String
		}
		if column.IsUnique.String == "YES" && !column.IsPrimary {
			query += " UNIQUE"
		}
		if column.IsNullable == "NO" && !column.IsPrimary {
			query += " NOT NULL"
		}
		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ")"
	return query
}

// convertTypeToMysql converts a Data type to a MySQL SQL Data type.
func convertTypeToMysql(dataType string) string {
	// Map column types to MySQL equivalents
	switch dataType {
	case "BOOL":
		return "TINYINT"
	case "BOOLEAN":
		return "TINYINT"
	case "CHARACTER VARYING":
		return "VARCHAR"
	case "FIXED":
		return "DECIMAL"
	case "FLOAT4":
		return "FLOAT"
	case "FLOAT8":
		return "DOUBLE"
	case "INT1":
		return "TINYINT"
	case "INT2":
		return "SMALLINT"
	case "INT3":
		return "MEDIUMINT"
	case "INT4":
		return "INT"
	case "INT8":
		return "BIGINT"
	case "LONG VARBINARY":
		return "MEDIUMBLOB"
	case "LONG VARCHAR", "LONG":
		return "MEDIUMTEXT"
	case "MIDDLEINT":
		return "MEDIUMINT"
	case "NUMERIC":
		return "DECIMAL"
	// Add more type conversions as needed
	default:
		return dataType
	}
}

// Create a new MySQL connection URL with the given configuration.
func dbURLMySQL(dbConfig *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?tls=%v&interpolateParams=true",
		dbConfig.Username,
		DB_PASSWORD,
		dbConfig.Host,
		dbConfig.Database,
		dbConfig.SSL,
	)
}
