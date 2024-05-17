package mssql

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thesaas-company/xray/types"
)

// setting up a mock db connection
// this function returns a mock database connection and a mock object
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred while creating a new mock database connection")
	}
	return db, mock
}

func TestSchema(t *testing.T) {
	db, mock := MockDB() // create a new mock database connection
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}() // close the connection when the function returns

	tableName := "user"                                                                  
	mockRows := sqlmock.NewRows([]string{"Field", "Type", "IsNullable", "ColumnDefault", "OrdinalPosition", "CharacterMaximumLength"}).AddRow("id", "int", "true", "", 1, nil)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(MSSQL_SCHEMA_QUERY, tableName))).WillReturnRows(mockRows) // set the expected return values for the query

	// we then create a new instance of our MySQL object and test the function
	m, err := NewMSSQL(db)
	if err != nil {
		t.Errorf("error initialising mysql: %s", err)
	}
	response, err := m.Schema(tableName) // call the Schema method
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}

	fmt.Printf("Table schema : %+v\n", response)

	// we make sure that all expectations were met, otherwise an error will be reported
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestTables(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	tableList := []string{"user", "product", "order"} // list of tables to be returned by the query
	schema := "test"
	// Retrieve the list of tables
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(MSSQL_TABLES_QUERY, schema))).WillReturnRows(rows) // set the expected return values for the query

	b, err := NewMSSQL(db) // create a new instance of our BigQuery object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := b.Tables(schema) // call the Tables method
	if err != nil {
		t.Errorf("error executing the query: %s", err)
	}

	fmt.Printf("Tables: %+v\n", tables)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


func TestGenerateCreateTablequery(t *testing.T) {
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	table := types.Table{
		Name: "user",
		Columns: []types.Column{
			{
				Name:         "id",
				Type:         "int",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    true,
				IsUnique:     sql.NullString{String: "YES", Valid: true},
			},
			{
				Name:         "name",
				Type:         "varchar(255)",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    false,
				IsUnique:     sql.NullString{String: "NO", Valid: true},
			},
			{
				Name:       "age",
				Type:       "int",
				IsNullable: "YES",
			},
		},
	}

	m := &MSSQL{Client: db}
	query := m.GenerateCreateTableQuery(table)

	fmt.Printf("Create table query: %s\n", query)

	expectedQuery := "CREATE TABLE user (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL, age INT)"
	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	// query to be executed
	query := `SELECT id,name FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John") // mock rows to be returned by the query

	mock.ExpectQuery(query).WillReturnRows(mockRows) // set the expected return values for the query

	m, err := NewMSSQL(db) // create a new instance of our MySQL object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	res, err := m.Execute(regexp.QuoteMeta(query)) // call the Execute method
	if err != nil {
		t.Errorf("error executing the query: %s", err)
	}

	var result types.QueryResult
	if err := json.Unmarshal(res, &result); err != nil {
		t.Errorf("error unmarshalling the result: %s", err)
	}

	fmt.Printf("Query result: %+v\n", result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
