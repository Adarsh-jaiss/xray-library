package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func TestSchema(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableName := "user"
	mockRows := sqlmock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra",}).AddRow("id", "int", "NO", "PRI", nil, "auto_increment")

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(SCHEMA_QUERY, tableName))).WillReturnRows(mockRows)

	// we then create a new instance of our MySQL object and test the function
	m, err := NewMySQL(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	response, err := m.Schema(tableName)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}

	fmt.Printf("Table schema : %+v\n", response)

	// we make sure that all expectations were met, otherwise an error will be reported
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT id,name FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John")

	mock.ExpectQuery(query).WillReturnRows(mockRows)

	m, err := NewMySQL(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	res, err := m.Execute(regexp.QuoteMeta(query))
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

func TestGetTableName(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableList := []string{"user", "product", "order"}

	// Retrieve the list of tables
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(MYSQL_TABLES_LIST_QUERY).WithArgs("test").WillReturnRows(rows)

	m, err := NewMySQL(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := m.Tables("test")
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	fmt.Printf("Table names: %+v\n", tables)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
