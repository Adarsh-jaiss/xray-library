package types

import (
	"database/sql"
)

// ISQL is an interface that defines the methods that a SQL database must implement.
type ISQL interface {
	Schema(string) (Table, error)          // Schema retrieves the schema for the specified table.
	Execute(string) ([]byte, error)        // Execute executes the given SQL query.
	Tables(string) ([]string, error)       // Tables retrieves the list of tables for the specified database.
	GenerateCreateTableQuery(Table) string // GenerateCreateTableQuery generates the CREATE TABLE query for the specified table.
}

// Table represents a database table.
type Table struct {
	Name        string   `json:"name"`         // Name is the name of the table.
	Dataset     string   `json:"dataset"`      // Dataset is the dataset of the bigquery table.
	Columns     []Column `json:"columns"`      // Columns are the columns in the table.
	ColumnCount int64    `json:"column_count"` // ColumnCount is the number of columns in the table.
	Description string   `json:"description"`  // Description is a description of the table.
	Metatags    []string `json:"metatags"`     // Metatags contains all column names.
}

// Column represents a column in a database table.
type Column struct {
	Name                   string         `json:"name"`                     // Name is the name of the column.
	Type                   string         `json:"type"`                     // Type is the data type of the column.
	AutoIncrement          bool           `json:"auto_increment"`           // AutoIncrement indicates whether the column is auto-incrementing.
	IsNullable             string         `json:"is_nullable"`              // IsNullable indicates whether the column can have null values.
	Key                    string         `json:"key"`                      // Key is the key type of the column.
	DefaultValue           sql.NullString `json:"default_value"`            // DefaultValue is the default value of the column.
	Extra                  string         `json:"extra"`                    // Extra contains additional information about the column.
	Description            string         `json:"description"`              // Description is a description of the column.
	Metatags               []string       `json:"metatags"`                 // Metatags contains the column name.
	Visibility             bool           `json:"visibility"`               // Visibility indicates whether the column is visible.
	IsIndex                bool           `json:"is_index"`                 // IsIndex indicates whether the column is an index.
	IsPrimary              bool           `json:"is_primary"`               // IsPrimary indicates whether the column is a primary key.
	ColumnDefault          sql.NullString `json:"column_default"`           // ColumnDefault is the default value of the column.
	CharacterMaximumLength sql.NullInt64  `json:"character_maximum_length"` // CharacterMaximumLength is the maximum length of the column.
	IsUpdatable            sql.NullString `json:"is_updatable"`             // IsUpdatable indicates whether the column is updatable.
	IsIdentity             sql.NullString `json:"is_identity"`              // IsIdentity indicates whether the column is an identity column.
	IsGenerated            sql.NullString `json:"is_generated"`             // IsGenerated indicates whether the column is a generated column.
	IsUnique               sql.NullString `json:"is_unique"`                // IsUnique indicates whether the column is unique.
	IsSystemColumn         sql.NullString `json:"is_system_column"`         // IsSystemColumn indicates whether the column is a system column.
	IsHidden               sql.NullString `json:"is_hidden"`                // IsHidden indicates whether the column is hidden.
	IsReadOnly             sql.NullString `json:"is_read_only"`             // IsReadOnly indicates whether the column is read-only.
	IsComputed             sql.NullString `json:"is_computed"`              // IsComputed indicates whether the column is computed.
	IsSparse               sql.NullString `json:"is_sparse"`                // IsSparse indicates whether the column is sparse.
	IsColumnSet            sql.NullString `json:"is_column_set"`            // IsColumnSet indicates whether the column is a column set.
	IsSelfReferencing      sql.NullString `json:"is_self_referencing"`      // IsSelfReferencing indicates whether the column is self-referencing.
	ScopeName              sql.NullString `json:"scope_name"`               // ScopeName is the scope name of the column.
	ScopeSchema            sql.NullString `json:"scope_schema"`             // ScopeSchema is the scope schema of the column.
	OrdinalPosition        sql.NullInt64  `json:"ordinal_position"`         // OrdinalPosition is the position of the column in the table.
	IdentitySeed           sql.NullInt64  `json:"identity_seed"`            // IdentitySeed is the seed value of the identity column.
	IdentityStep           sql.NullInt64  `json:"identity_step"`            // IdentityStep is the step value of the identity column.
}

// TableResponse is the struct that holds the response from the Redshift Tables query.
type TableResponse struct {
	TableName    string         `json:"table_name"`    // TableName is the name of the table.
	DatabaseName string         `json:"database_name"` // DatabaseName is the name of the database.
	SchemaName   string         `json:"schema_name"`   // SchemaName is the name of the schema.
	TableType    string         `json:"table_type"`    // TableType is the type of the table.
	TableAcl     sql.NullString `json:"table_acl"`     // TableAcl is the access control list of the table.
	Remarks      sql.NullString `json:"remarks"`       // Remarks is the remarks of the table.
}

// QueryResult represents the result of a database query.
type QueryResult struct {
	Columns []string        `json:"columns"` // Columns are the names of the columns in the result.
	Rows    [][]interface{} `json:"rows"`    // Rows are the rows in the result.
	Time    int64           `json:"time"`    // Time is the time it took to execute the query.
	Error   string          `json:"error"`   // Error is any error that occurred while executing the query.
}

// DbType represents a type of SQL database.
type DbType int

// These constants represent the supported types of SQL databases.
const (
	MySQL DbType = iota + 1
	Postgres
	Snowflake
	BigQuery
	Redshift
	MSSQL
)

// String returns the string representation of the DbType.
func (w DbType) String() string {
	return [...]string{"mysql", "postgres", "snowflake", "bigquery", "redshift", "mssql"}[w-1]
}

// Index returns the integer representation of the DbType.
func (w DbType) Index() int {
	return int(w)
}
