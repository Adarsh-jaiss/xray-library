package snowflake

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sf "github.com/snowflakedb/gosnowflake"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

type Snowflake struct {
	Client *sql.DB
	Config *config.Config
}

const (
	SNOWFLAKE_PASSWORD = "root"
	// SNOWFLAKE_SCHEMA_QUERY      = "DESCRIBE %s"
	SNOWFLAKE_TABLES_LIST_QUERY = "SHOW TERSE TABLES"
	SNOWFLAKE_SCHEMA_QUERY      = `
	SELECT 
		COLUMN_NAME, 
		DATA_TYPE, 
		IS_NULLABLE, 
		COLUMN_DEFAULT, 
		CHARACTER_MAXIMUM_LENGTH, 
		IS_UPDATABLE, 
		IS_IDENTITY, 
		IS_GENERATED, 
		IS_UNIQUE, 
		IS_SYSTEM_COLUMN, 
		IS_HIDDEN, 
		IS_READ_ONLY, 
		IS_COMPUTED, 
		IS_SPARSE, 
		IS_COLUMN_SET, 
		IS_SELF_REFERENCING, 
		SCOPE_NAME, 
		SCOPE_SCHEMA, 
		ORDINAL_POSITION 
	FROM INFORMATION_SCHEMA.COLUMNS 
	WHERE TABLE_NAME = ?)
	`
)

// The NewSnowflake function is responsible for creating a new Snowflake object with an initialized database client and configuration.
func NewSnowflake(dbClient *sql.DB) (types.ISQL, error) {
	return &Snowflake{
		Client: dbClient,
		Config: &config.Config{},
	}, nil
}

func NewSnowflakeWithConfig(config *config.Config) (types.ISQL, error) {
	dsn, err := sf.DSN(&sf.Config{
		Account:   config.Account,
		User:      config.Username,
		Password:  SNOWFLAKE_PASSWORD,
		Database:  config.DatabaseName,
		Warehouse: config.Warehouse,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating snowflake DSN: %v", err)
	}

	dbType := types.Snowflake
	db, err := sql.Open(dbType.String(), dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to snowflake database: %v", err)
	}

	return &Snowflake{
		Client: db,
		Config: config,
	}, nil

}

func (s *Snowflake) Schema(table string) (types.Table, error) {
	var res types.Table

	rows, err := s.Client.Query(SNOWFLAKE_SCHEMA_QUERY, table)
	if err != nil {
		return res, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&column.IsNullable,
			&column.ColumnDefault,
			&column.CharacterMaximumLength,
			&column.IsUpdatable,
			&column.IsIdentity,
			&column.IsGenerated,
			&column.IsUnique,
			&column.IsSystemColumn,
			&column.IsHidden,
			&column.IsReadOnly,
			&column.IsComputed,
			&column.IsSparse,
			&column.IsColumnSet,
			&column.IsSelfReferencing,
			&column.ScopeName,
			&column.ScopeSchema,
			&column.OrdinalPosition); err != nil {
			return res, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Description = ""      // default description
		column.Metatags = []string{} // default metatags as an empty string slice
		column.Metatags = append(column.Metatags, column.Name)
		column.Visibility = true // default visibility
		columns = append(columns, column)
	}

	// checking for erros from iterating over the rows
	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil
}

// Every table in Snowflake lives "inside" a schema. Every schema lives "inside" a database. It's a hierarchical system.
func (s *Snowflake) Tables(SchemaName string) ([]string, error) {
	query := fmt.Sprintf("USE WAREHOUSE %s", s.Config.Warehouse)
	rows, err := s.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning database: %v", err)
		}
		tables = append(tables, table)
	}

	// checking for errors in iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows:%v", err)
	}

	return tables, nil
}

func (s *Snowflake) Execute(query string) ([]byte, error) {
	rows, err := s.Client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing sql statement: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		results = append(results, values)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Convert the result to JSON
	queryResult := types.QueryResult{
		Columns: columns,
		Rows:    results,
	}

	jsonData, err := json.Marshal(queryResult)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %v", err)
	}

	return jsonData, nil
}
