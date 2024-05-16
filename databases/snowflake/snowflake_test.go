package snowflake

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thesaas-company/xray/types"
)

// MockDB is a mock implementation of the Snowflake struct.
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred while creating a new mock database connection")
	}

	return db, mock
}

// TestSchema is a unit test function that tests the Schema method of the Snowflake struct.
// It creates a mock instance of Snowflake, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestSchema(t *testing.T) {
	db, mock := MockDB() // create a new mock database connection
	defer db.Close()

	table_name := "user"

	// mock rows to be returned by the query
	columns := []string{"name", "type"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int").AddRow("name", "varchar")
	// set the expected return values for the query
	mock.ExpectQuery(regexp.QuoteMeta(SNOWFLAKE_SCHEMA_QUERY)).WithArgs(table_name).WillReturnRows(mockRows)

	s, err := NewSnowflake(db) // create a new instance of our Snowflake object
	if err != nil {
		t.Errorf("error initialising snowflake: %s", err)
	}

	res, err := s.Schema(table_name) // call the Schema method
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	fmt.Printf("Table schema %+v\n", res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}

}

// TestExecute is a unit test function that tests the Execute method of the Snowflake struct.
// It creates a mock instance of Snowflake, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestExecute(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT id, name FROM "user"`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan") // mock rows to be returned by the query

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows) // set the expected return values for the query

	p, err := NewSnowflake(db) // create a new instance of our Snowflake object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	res, err := p.Execute(query) // call the Execute method
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

// TestTables is a unit test function that tests the Tables method of the Snowflake struct.
// It creates a mock instance of Snowflake, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestTables(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer db.Close()

	tableList := []string{"user", "product", "order"}
	Warehouse := "datasherlock"
	// set the expected return values for the query
	mock.ExpectQuery("USE WAREHOUSE ").WithArgs(Warehouse).WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(""))

	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(SNOWFLAKE_TABLES_LIST_QUERY)).WillReturnRows(rows)

	s, err := NewSnowflake(db) // create a new instance of our Snowflake object
	if err != nil {
		t.Fatalf("error initializing snowflake: %s", err)
	}

	query := fmt.Sprintf("USE WAREHOUSE %s", Warehouse)
	_, err = s.Tables(query) // call the Tables method
	if err != nil {
		return
	}

	tables, err := s.Tables("test") // Database name isn't used in the query, so you can pass any value here
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	expected := tableList // Using the same list as returned by the mock
	if !reflect.DeepEqual(tables, expected) {
		t.Errorf("expected: %v, got: %v", expected, tables)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenerateCreateTableQuery(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	table := types.Table{
		Name: "user",
		Columns: []types.Column{
			{
				Name:          "id",
				Type:          "int",
				AutoIncrement: true,
				IsNullable:    "NO",
				DefaultValue:  sql.NullString{String: "", Valid: false},
				IsPrimary:     true,
				IsUnique:      sql.NullString{String: "YES", Valid: true},
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

	s := &Snowflake{}
	query := s.GenerateCreateTableQuery(table)

	expectedQuery := "CREATE TABLE user (id INT AUTOINCREMENT PRIMARY KEY UNIQUE, name VARCHAR(255) NOT NULL, age INT);"
	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
