package bigquery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	// "github.com/thesaas-company/xray/config"
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

// TestSchema is a unit test function that tests the Schema method of the BigQuery struct.
// It creates a mock instance of BigQuery, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestSchema(t *testing.T) {
	db, mock := MockDB() // create a new mock database connection
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}() // close the connection when the function returns

	tableName := "user"                                                                                                                               // table name to be used in the test
	mockRows := sqlmock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).AddRow("id", "int", "NO", "PRI", nil, "auto_increment") // mock rows to be returned by the query

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(BigQuery_SCHEMA_QUERY, tableName))).WillReturnRows(mockRows) // set the expected return values for the query

	// we then create a new instance of our BigQuery object and test the function
	m, err := NewBigQuery(db)
	if err != nil {
		t.Errorf("error initialising bigquery: %s", err)
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

// TestTables is a unit test function that tests the Tables method of the BigQuery struct.
// It creates a mock instance of BigQuery, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestTables(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	tableList := []string{"user", "product", "order"} // list of tables to be returned by the query
	schema := "test"
	// Retrieve the list of tables
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(BigQuery_TABLES_QUERY, schema))).WillReturnRows(rows) // set the expected return values for the query

	b, err := NewBigQuery(db) // create a new instance of our BigQuery object
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


// TestExecute is a unit test function that tests the Execute method of the BigQuery struct.
// It creates a mock instance of BigQuery, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestExecute(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// query to be executed
	query := "SELECT * FROM `my_project.test.test.my_table`"

	// Prepare the rows that will be returned by the query
	mockRows := sqlmock.NewRows([]string{"column1", "column2"}).
		AddRow("value1", "value2").
		AddRow("value3", "value4")

	// Set the expected return values for the query
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	// create a new instance of our BigQuery object
	b, err := NewBigQuery(db)
	if err != nil {
		t.Errorf("error initializing bigquery: %s", err)
	}

	// call the Execute method
	jsonData, err := b.Execute(query)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}

	// Unmarshal the JSON data
	var result types.BigQueryResult
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Errorf("error unmarshaling json: %s", err)
	}

	// Check if the returned columns match the expected columns
	expectedColumns := []string{"column1", "column2"}
	if !reflect.DeepEqual(result.Columns, expectedColumns) {
		t.Errorf("expected %v, got %v", expectedColumns, result.Columns)
	}

	// Check if the returned rows match the expected rows
	expectedRows := []map[string]interface{}{
		{"column1": "value1", "column2": "value2"},
		{"column1": "value3", "column2": "value4"},
	}
	if !reflect.DeepEqual(result.Rows, expectedRows) {
		t.Errorf("expected %v, got %v", expectedRows, result.Rows)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGenerateCreateTablequeryis a unit test function that tests the TestGenerateCreateTablequery method of the BigQuery struct.
// It creates a mock instance of BigQuery, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestGenerateCreateTablequery(t *testing.T) {
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	table := types.Table{
		Dataset: "test_dataset",
		Name:    "table_name",
		Columns: []types.Column{
			{Name: "id", Type: "INT64", IsPrimary: true},
			{Name: "name", Type: "STRING"},
			{Name: "created_at", Type: "TIMESTAMP"},
		},
	}

	b, err := NewBigQuery(db)
	if err != nil {
		t.Errorf("error initializing bigquery: %s", err)
	}
	query := b.GenerateCreateTableQuery(table)
	fmt.Println(query)

	expectedQuery := "CREATE TABLE test_dataset.table_name (id INT64, name STRING, created_at TIMESTAMP);"
	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
