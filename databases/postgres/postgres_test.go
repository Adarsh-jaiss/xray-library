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

	table_name := "user"

	columns := []string{"name", "type", "IsNullable", "key", "Description", "Extra", "IsPrimary", "IsIndex"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", "No", "PRIMARY", "This is the primary key of the table to identify users", "auto_increment", "true", "true")

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(POSTGRES_SCHEMA_QUERY, table_name))).WillReturnRows(mockRows)

	m, err := NewPostgres(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	response, err := m.Schema(table_name)
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	fmt.Printf("Table schema: %+v\n", response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT id, name FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	p, err := NewPostgres(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	res, err := p.Execute(query)
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

	tableList := []string{"user", "Credit", "Debit"}
	DatabaseName := "test"
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(POSTGRES_TABLE_LIST_QUERY, DatabaseName))).WillReturnRows(rows)

	p, err := NewPostgres(db)
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := p.Tables(DatabaseName)
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
