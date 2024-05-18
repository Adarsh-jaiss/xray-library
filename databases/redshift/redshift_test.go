package redshift

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
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	table_name := "user" // table name to be used in the test

	// mock rows to be returned by the query
	columns := []string{"column", "type", "encoding", "distkey", "sortkey", "notnull"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", "utf8", true, 1, true)
	// set the expected return values for the query
	expectedQuery := fmt.Sprintf(Redshift_Schema_query, "public", table_name)
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(mockRows)

	// we then create a new instance of our Redshift object and test the function
	r, err := NewRedshift(db)
	if err != nil {
		t.Errorf("error initialising redshift: %s", err)
	}
	response, err := r.Schema(table_name) // call the Schema method
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
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//	we then create a new instance of our Postgres object and test the function
	query := `SELECT id, name FROM "user"`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan") // mock rows to be returned by the query

	// set the expected return values for the query
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	// create a new instance of our Postgres object
	r, err := NewRedshift(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	res, err := r.Execute(query) // call the Execute method
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
func TestTables(t *testing.T) {
	db, mock := MockDB() // create a new mock database connection
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	tableList := []string{"user", "Credit", "Debit"} // list of tables to be used in the test
	DatabaseName := "test"
	rows := sqlmock.NewRows([]string{"DatabaseName", "SchemaName", "TableName", "TableType", "TableAcl", "Remarks"}).
		AddRow(DatabaseName, "public", tableList[0], "BASE TABLE", nil, nil).
		AddRow(DatabaseName, "public", tableList[1], "BASE TABLE", nil, nil).
		AddRow(DatabaseName, "public", tableList[2], "BASE TABLE", nil, nil)
	expectedQuery := fmt.Sprintf(Redshift_Tables_query, DatabaseName)
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)

	r, err := NewRedshift(db) // create a new instance of our Postgres object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := r.Tables(DatabaseName) // call the Tables method
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

// TestGenerateCreateTableQuery is a unit test function that tests the GenerateCreateTableQuery method of the Redshift struct.
// It creates a mock instance of Redshift, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestGenerateCreateTableQuery(t *testing.T) {
	db, mock := MockDB()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

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

	r := Redshift{Client: db}
	query := r.GenerateCreateTableQuery(table)

	fmt.Printf("Create table query: %v\n", query)

	expectedQuery := `CREATE TABLE ..user (id INT PRIMARY KEY IDENTITY({0 false}, {0 false}) NOT NULL, name VARCHAR(255) NOT NULL, age INT);`
	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

}
