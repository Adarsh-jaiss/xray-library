package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thesaas-company/xray/config"
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

// TestSchema is a unit test function that tests the Schema method of the MySQL struct.
// It creates a mock instance of MySQL, sets the expected return values, and calls the method under test.
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

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(SCHEMA_QUERY, tableName))).WillReturnRows(mockRows) // set the expected return values for the query

	// we then create a new instance of our MySQL object and test the function
	m, err := NewMySQL(db)
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

// TestExecute is a unit test function that tests the Execute method of the MySQL struct.
// It creates a mock instance of MySQL, sets the expected return values, and calls the method under test.
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
	query := `SELECT id,name FROM user`
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John") // mock rows to be returned by the query

	mock.ExpectQuery(query).WillReturnRows(mockRows) // set the expected return values for the query

	m, err := NewMySQL(db) // create a new instance of our MySQL object
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

// TestGetTableName is a unit test function that tests the Tables method of the MySQL struct.
// It creates a mock instance of MySQL, sets the expected return values, and calls the method under test.
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

	// Retrieve the list of tables
	rows := sqlmock.NewRows([]string{"table_name"}).
		AddRow(tableList[0]).
		AddRow(tableList[1]).
		AddRow(tableList[2])
	mock.ExpectQuery(MYSQL_TABLES_LIST_QUERY).WithArgs("test").WillReturnRows(rows) // set the expected return values for the query

	m, err := NewMySQL(db) // create a new instance of our MySQL object
	if err != nil {
		t.Errorf("error executing query: %s", err)
	}
	tables, err := m.Tables("test") // call the Tables method
	if err != nil {
		t.Errorf("error retrieving table names: %s", err)
	}

	fmt.Printf("Table names: %+v\n", tables)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGenerateCreateTablequery is a unit test function that tests the GenerateCreateTableQuery method of the MySQL.
// It creates a mock instance of MySQL, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestGenerateCreateTablequery(t *testing.T) {
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

	m := &MySQL{Client: db}
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

func TestNewMySQL(t *testing.T) {
	type args struct {
		dbClient *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    types.ISQL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMySQL(tt.args.dbClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMySQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMySQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMySQLWithConfig(t *testing.T) {
	type args struct {
		dbConfig *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    types.ISQL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMySQLWithConfig(tt.args.dbConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMySQLWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMySQLWithConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMySQL_Schema(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name    string
		m       *MySQL
		args    args
		want    types.Table
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Schema(tt.args.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("MySQL.Schema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MySQL.Schema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMySQL_Execute(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		m       *MySQL
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Execute(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("MySQL.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MySQL.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMySQL_Tables(t *testing.T) {
	type args struct {
		databaseName string
	}
	tests := []struct {
		name    string
		m       *MySQL
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Tables(tt.args.databaseName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MySQL.Tables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MySQL.Tables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMySQL_GenerateCreateTableQuery(t *testing.T) {
	type args struct {
		table types.Table
	}
	tests := []struct {
		name string
		m    *MySQL
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GenerateCreateTableQuery(tt.args.table); got != tt.want {
				t.Errorf("MySQL.GenerateCreateTableQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertTypeToMysql(t *testing.T) {
	type args struct {
		dataType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTypeToMysql(tt.args.dataType); got != tt.want {
				t.Errorf("convertTypeToMysql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbURLMySQL(t *testing.T) {
	type args struct {
		dbConfig *config.Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dbURLMySQL(tt.args.dbConfig); got != tt.want {
				t.Errorf("dbURLMySQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
