package redshift

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thesaas-company/xray/types"
)


type MockRedshift struct {
	mock.Mock
}

func (m *MockRedshift) Schema(table string) (types.Table, error) {
	args := m.Called(table)
	return args.Get(0).(types.Table), args.Error(1)
}

func (m *MockRedshift) Execute(query string) ([]byte, error) {
	args := m.Called(query)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRedshift) Tables(database string) ([]string, error) {
	args := m.Called(database)
	return args.Get(0).([]string), args.Error(1)
}

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

func TestRedshift_Execute(t *testing.T) {
	mockRedshift := new(MockRedshift)

	expectedResult := []byte("result")
	mockRedshift.On("Execute", "query").Return(expectedResult, nil)

	actualResult, err := mockRedshift.Execute("query")

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestRedshift_Tables(t *testing.T) {
	mockRedshift := new(MockRedshift)

	expectedTables := []string{"table1", "table2"}
	mockRedshift.On("Tables", "database").Return(expectedTables, nil)

	actualTables, err := mockRedshift.Tables("database")

	assert.NoError(t, err)
	assert.Equal(t, expectedTables, actualTables)
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
