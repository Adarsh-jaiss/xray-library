package redshift

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thesaas-company/xray/types"
)

// MockRedshift is a mock implementation of the Redshift interface.
type MockRedshift struct {
	mock.Mock
}

// Schema mocks the Schema method.
func (m *MockRedshift) Schema(table string) (types.Table, error) {
	args := m.Called(table)
	return args.Get(0).(types.Table), args.Error(1)
}

// Execute mocks the Execute method.
func (m *MockRedshift) Execute(query string) ([]byte, error) {
	args := m.Called(query)
	return args.Get(0).([]byte), args.Error(1)
}

// Tables mocks the Tables method.
func (m *MockRedshift) Tables(database string) ([]string, error) {
	args := m.Called(database)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedshift) GenerateCreateTableQuery(table types.Table) string {
	args := m.Called(table)
	return args.Get(0).(string)
}

// Redshidft schema
func TestRedshift_Schema(t *testing.T) {
	mockRedshift := new(MockRedshift)

	expectedSchema := types.Table{
		Name:        "table_name",
		Description: "table_description",
		Columns: []types.Column{
			{
				Name: "column1",
				Type: "string",
			},
			{
				Name: "column2",
				Type: "int",
			},
		},
	}
	mockRedshift.On("Schema", "table_name").Return(expectedSchema, nil)

	actualSchema, err := mockRedshift.Schema("table_name")

	assert.NoError(t, err)
	assert.Equal(t, expectedSchema, actualSchema)
}

// Redhshift schema is a unit test function that tests the Schema method of the Redhsift struct.
// It creates a mock instance of Redshift, sets the expected return values, and calls the method under test.
// It then asserts the expected return values and checks if the method was called with the correct arguments.
func TestRedshift_Execute(t *testing.T) {
	mockRedshift := new(MockRedshift)

	expectedResult := []byte("result")
	mockRedshift.On("Execute", "query").Return(expectedResult, nil)

	actualResult, err := mockRedshift.Execute("query")

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

// Redshift schema is a unit test function that tests the Tables method of the Redshift struct.
func TestRedshift_Tables(t *testing.T) {
	mockRedshift := new(MockRedshift)

	expectedTables := []string{"table1", "table2"}
	mockRedshift.On("Tables", "database").Return(expectedTables, nil)

	actualTables, err := mockRedshift.Tables("database")

	assert.NoError(t, err)
	assert.Equal(t, expectedTables, actualTables)
}

func TestGenerateCreateTableQuery(t *testing.T) {
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

	// Define the expected query
	expectedQuery := "CREATE TABLE amazon.customer.user (id INT PRIMARY KEY IDENTITY(1,1) UNIQUE NOT NULL, name VARCHAR(255) NOT NULL, age INT);"

	// Create a mock Redshift instance
	mockRedshift := new(MockRedshift)
	mockRedshift.On("GenerateCreateTableQuery", table).Return(expectedQuery)

	// Generate the create table query
	actualQuery := mockRedshift.GenerateCreateTableQuery(table)

	// Compare the generated query with the expected query
	assert.Equal(t, expectedQuery, actualQuery)
}

func TestRedshift_Schema_Error(t *testing.T) {
	mockRedshift := new(MockRedshift)

	mockRedshift.On("Schema", "table_name").Return(types.Table{}, errors.New("error"))

	_, err := mockRedshift.Schema("table_name")

	assert.Error(t, err)
}

func TestRedshift_Execute_Error(t *testing.T) {
	mockRedshift := new(MockRedshift)

	mockRedshift.On("Execute", "query").Return([]byte{}, errors.New("error"))

	_, err := mockRedshift.Execute("query")

	assert.Error(t, err)
}

func TestRedshift_Tables_Error(t *testing.T) {
	mockRedshift := new(MockRedshift)

	mockRedshift.On("Tables", "database").Return([]string{}, errors.New("error"))

	_, err := mockRedshift.Tables("database")

	assert.Error(t, err)
}

func TestRedshift_RedshiftAPIService(t *testing.T) {
	query := "query"

	mockRedshift := new(MockRedshift)
	expectedResult := []byte("result")
	mockRedshift.On("Execute", query).Return(expectedResult, nil)
	actualResult, err := mockRedshift.Execute(query)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}
