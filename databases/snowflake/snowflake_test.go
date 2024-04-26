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

func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occured while creating a new mock database connection")
	}

	return db, mock
}

func TestSchema(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	table_name := "user"

	columns := []string{"name", "type", "IsNullable", "DefaultValue", "IsUpdatable", "IsIdenity", "IsGenerated", "IsUnique", "IsSystemColumn", "IsHidden", "IsReadOnly", "IsComputed", "IsSparse", "IsColumnSet", "IsSelfReplacing", "ScopeName", "ScopeSchema", "OrdinalPosition"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", true, 1, true, false, true, true, false, false, true, false, true, false, true, "scope1", "schema1", 1)

	mock.ExpectQuery(regexp.QuoteMeta(SNOWFLAKE_SCHEMA_QUERY)).WithArgs(table_name).WillReturnRows(mockRows)

	s, err := NewSnowflake(db)
	if err != nil {
		t.Errorf("error initialising snowflake: %s", err)
	}

	res, err := s.Schema(table_name)
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	fmt.Printf("Table schema %+v\n", res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}

}

func TestExecute(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	query := `SELECT id, name FROM "user"`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Rohan")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	p, err := NewSnowflake(db)
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

func TestTables(t *testing.T) {
	db, mock := MockDB()
	defer db.Close()

	tableList := []string{"user", "product", "order"}
	Warehouse := "datasherlock"
	mock.ExpectQuery("USE WAREHOUSE ").WithArgs(Warehouse).WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(""))

	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(regexp.QuoteMeta(SNOWFLAKE_TABLES_LIST_QUERY)).WillReturnRows(rows)

	s, err := NewSnowflake(db)
	if err != nil {
		t.Fatalf("error initializing snowflake: %s", err)
	}

	query := fmt.Sprintf("USE WAREHOUSE %s", Warehouse)
	_, err = s.Tables(query)
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
