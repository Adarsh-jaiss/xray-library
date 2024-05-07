package bigquery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thesaas-company/xray/types"
)

// MockBigQuery is a mock implementation of the BigQuery struct.
type MockBigQuery struct {
	mock.Mock
}

// Schema is a mock implementation of the Schema method.
func (m *MockBigQuery) Schema(table string) (types.Table, error) {
	args := m.Called(table)
	return args.Get(0).(types.Table), args.Error(1)
}

// Execute is a mock implementation of the Execute method.
func (m *MockBigQuery) Execute(query string) ([]byte, error) {
	args := m.Called(query)
	return args.Get(0).([]byte), args.Error(1)
}

// Tables is a mock implementation of the Tables method.
func (m *MockBigQuery) Tables(dataset string) ([]string, error) {
	args := m.Called(dataset)
	return args.Get(0).([]string), args.Error(1)
}

func TestBigQuery_Schema(t *testing.T) {
	// Create a new instance of the mock
	mockBigQuery := new(MockBigQuery)

	// Set the expected return values
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
	mockBigQuery.On("Schema", "table_name").Return(expectedSchema, nil)

	// Call the method under test
	actualSchema, err := mockBigQuery.Schema("table_name")

	// Assert the expected return values
	assert.NoError(t, err)
	assert.Equal(t, expectedSchema, actualSchema)

	// Assert that the method was called with the correct arguments
	mockBigQuery.AssertCalled(t, "Schema", "table_name")
	fmt.Println(expectedSchema, actualSchema)
}

func TestBigQuery_Execute(t *testing.T) {
	// Create a new instance of the mock
	mockBigQuery := new(MockBigQuery)

	// Set the expected return values
	expectedResult := []byte(`{"result": "success"}`)
	mockBigQuery.On("Execute", "SELECT * FROM table").Return(expectedResult, nil)

	// Call the method under test
	actualResult, err := mockBigQuery.Execute("SELECT * FROM table")

	// Assert the expected return values
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)

	// Assert that the method was called with the correct arguments
	mockBigQuery.AssertCalled(t, "Execute", "SELECT * FROM table")
	fmt.Println(expectedResult, actualResult)
}

func TestBigQuery_Tables(t *testing.T) {
	// Create a new instance of the mock
	mockBigQuery := new(MockBigQuery)

	// Set the expected return values
	expectedTables := []string{"table1", "table2"}
	mockBigQuery.On("Tables", "dataset").Return(expectedTables, nil)

	// Call the method under test
	actualTables, err := mockBigQuery.Tables("dataset")

	// Assert the expected return values
	assert.NoError(t, err)
	assert.Equal(t, expectedTables, actualTables)

	// Assert that the method was called with the correct arguments
	mockBigQuery.AssertCalled(t, "Tables", "dataset")
	fmt.Println(expectedTables, actualTables)
}
