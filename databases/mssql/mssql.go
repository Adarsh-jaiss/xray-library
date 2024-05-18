package mssql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

var DB_PASSWORD = "DB_PASSWORD"

const (
	MSSQL_SCHEMA_QUERY = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT, ORDINAL_POSITION, CHARACTER_MAXIMUM_LENGTH FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s'"
	MSSQL_TABLES_QUERY = "USE %s; SELECT table_name FROM INFORMATION_SCHEMA.TABLES;"
)

type MSSQL struct {
	Client *sql.DB
	Config *config.Config
}

func NewMSSQL(client *sql.DB) (types.ISQL, error) {
	return &MSSQL{
		Client: client,
		Config: &config.Config{},
	}, nil
}

func NewMSSQLFromConfig(config *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 { // added mysql to be more verbose about the db type
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}

	DB_PASSWORD = os.Getenv(DB_PASSWORD)
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s", config.Host, config.Username, DB_PASSWORD, config.Port)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}

	return &MSSQL{
		Client: conn,
		Config: config,
	}, nil
}

func (m *MSSQL) Schema(table string) (types.Table, error) {
	query := fmt.Sprintf(MSSQL_SCHEMA_QUERY, table)
	rows, err := m.Client.Query(query)
	if err != nil {
		return types.Table{}, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	var columns []types.Column
	for rows.Next() {
		var col types.Column
		if err := rows.Scan(
			&col.Name,
			&col.Type,
			&col.IsNullable,
			&col.ColumnDefault,
			&col.OrdinalPosition,
			&col.CharacterMaximumLength,
		); err != nil {
			return types.Table{}, fmt.Errorf("error scanning rows : %v", err)
		}
		col.Description = ""      // default description
		col.Metatags = []string{} // default metatags as an empty string slice
		col.Metatags = append(col.Metatags, col.Name)
		col.Visibility = true // default visibility
		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return types.Table{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Metatags:    []string{},
	}, nil
}

func (m *MSSQL) Tables(databaseName string) ([]string, error) {
	query := fmt.Sprintf(MSSQL_TABLES_QUERY, databaseName)
	rows, err := m.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing the sql statement: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("error closing the rows: %v", err)
		}
	}()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning the database :%v", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows :%v", err)
	}

	return tables, nil

}

func (m *MSSQL) Execute(query string) ([]byte, error) {
	rows, err := m.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing the sql statement %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns : %v", err)
	}

	var results [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning rows:%v", err)
		}

		results = append(results, values)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows : %v", err)
	}

	queryResult := types.QueryResult{
		Columns: columns,
		Rows:    results,
	}

	jsonData, err := json.Marshal(queryResult)
	if err != nil {
		return nil, fmt.Errorf("error marshalling json: %v", err)
	}

	return jsonData, nil
}

func (m *MSSQL) GenerateCreateTableQuery(table types.Table) string {
	query := "CREATE TABLE [" + table.Name + "] ("
	pk := ""
	unique := ""
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += "[" + column.Name + "] " + colType
		if column.AutoIncrement {
			query += " IDENTITY(1,1)"
		}
		if column.IsPrimary {
			pk = " PRIMARY KEY ([" + column.Name + "])"
		}
		if column.DefaultValue.Valid {
			query += " DEFAULT (" + column.DefaultValue.String + ")"
		}
		if column.IsUnique.String == "YES" && !column.IsPrimary {
			unique = ", UNIQUE ([" + column.Name + "])"
		}
		if column.IsNullable == "NO" && !column.IsPrimary {
			query += " NOT NULL"
		}
		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += pk + unique + ")"
	return query
}
