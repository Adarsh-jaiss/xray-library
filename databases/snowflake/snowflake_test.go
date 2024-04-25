package snowflake

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	columns := []string{"name", "type", "IsNullable","DefaultValue", "characterMaximumLenghth", "IsUpdatable", "IsIdenity","IsGenerated", "IsUnique", "IsSystemColumn", "IsHidden", "IsReadOnly", "IsComputed", "IsSparse","IsColumnSet","IsSelfReplacing","ScopeName","ScopeSchema","OrdinalPosition"}
	mockRows := sqlmock.NewRows(columns).AddRow("id", "int", true, 1, "", true, false, true, true, false, false, true, false, true, false, true, "scope1", "schema1", 1)

	mock.ExpectQuery(regexp.QuoteMeta(SNOWFLAKE_SCHEMA_QUERY)).WithArgs(table_name).WillReturnRows(mockRows)

	s, err := NewSnowflake(db)
	if err != nil {
		t.Errorf("error initialising snowflake: %s",err)
	}

	res, err := s.Schema(table_name)
	if err != nil {
		t.Errorf("error executing query : %v", err)
	}

	fmt.Printf("Table schema %+v\n",res)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}

}
