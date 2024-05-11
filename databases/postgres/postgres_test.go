// This is unit testing of postgres using a mock DB

package postgres

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

// setting up a mock db connection
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred while creating a new mock database connection")
	}

	return db, mock
}

// TestSchema is a unit test function that tests the Schema method of the Postgres struct.
// It creates a mock instance of Postgres, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestSchema(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer db.Close()

	table_name := "user" // table name to be used in the test

	// mock rows to be returned by the query
	columns := []string{"name", "type", "IsNullable", "DefaultValue", "CharacterMaximumLength", "OrdinalPosition", "Visibility", "IsPrimary", "IsUpdatable"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", "No", "", 0, 1, true, true, true)
	// set the expected return values for the query
	mock.ExpectQuery(regexp.QuoteMeta(POSTGRES_SCHEMA_QUERY)).WithArgs(table_name).WillReturnRows(mockRows)

	// we then create a new instance of our Postgres object and test the function
	m, err := NewPostgres(db)
	if err != nil {
		t.Errorf("error initialising postgres: %s", err)
	}
	response, err := m.Schema(table_name) // call the Schema method
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	fmt.Printf("Table schema: %+v\n", response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}
}

// TestExecute is a unit test function that tests the Execute method of the Postgres struct.
// It creates a mock instance of Postgres, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestExecute(t *testing.T) {
	// create a new mock database connection
	db, mock := MockDB()
	defer db.Close()

	//	we then create a new instance of our Postgres object and test the function
	query := `SELECT id, name FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan") // mock rows to be returned by the query

	// set the expected return values for the query
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	// create a new instance of our Postgres object
	p, err := NewPostgres(db)
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

// TestGetTableName is a unit test function that tests the Tables method of the Postgres struct.
// It creates a mock instance of Postgres, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestGetTableName(t *testing.T) {
	db, mock := MockDB() // create a new mock database connection
	defer db.Close()

	tableList := []string{"user", "Credit", "Debit"} // list of tables to be used in the test
	DatabaseName := "test"
	rows := sqlmock.NewRows([]string{"table_name"}). // mock rows to be returned by the query
								AddRow(tableList[0]).
								AddRow(tableList[1]).
								AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(POSTGRES_TABLE_LIST_QUERY)).WithArgs(DatabaseName).WillReturnRows(rows)

	p, err := NewPostgres(db) // create a new instance of our Postgres object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := p.Tables(DatabaseName) // call the Tables method
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	expected := []string{"user", "Credit", "Debit"}
	if !reflect.DeepEqual(tables, expected) {
		t.Errorf("expected: %v, got: %v", expected, tables)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenerateCreateTablequery(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	table := types.Table{
		Name: "user",
		Columns: []types.Column{
			{
				Name:         "id",
				Type:         "SERIAL",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    true,
				IsUnique:     sql.NullString{String: "YES", Valid: true},
			},
			{
				Name:         "name",
				Type:         "VARCHAR(255)",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    false,
				IsUnique:     sql.NullString{String: "NO", Valid: true},
			},
			{
				Name:       "age",
				Type:       "INTEGER",
				IsNullable: "YES",
			},
		},
	}

	p := Postgres{Client: db}
	query := p.GenerateCreateTableQuery(table)

	fmt.Printf("Create table query: %v\n", query)

	expectedQuery := `CREATE TABLE "user" (id SERIAL PRIMARY KEY UNIQUE NOT NULL, name VARCHAR(255) NOT NULL, age INTEGER);`
	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}	
}
