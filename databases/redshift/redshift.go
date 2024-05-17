package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// DB_PASSWORD is the environment variable that holds the database password.
var DB_PASSWORD = "DB_PASSWORD"

const (
	Redshift_Schema_query = "SHOW COLUMNS FROM TABLE %s.%s.%s;"
	Redshift_Tables_query = "SELECT * FROM svv_all_tables WHERE database_name = '%s';"
)

type Redshift struct {
	Client *sql.DB
	Config config.Config
}

func NewRedshift(client *sql.DB) (types.ISQL, error) {
	return &Redshift{
		Client: client,
		Config: config.Config{},
	}, nil
}

func NewRedshiftWithConfig(cfg *config.Config) (types.ISQL, error) {
	if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 {
		return nil, fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
	}
	DB_PASSWORD = os.Getenv(DB_PASSWORD)

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, DB_PASSWORD, cfg.DatabaseName, cfg.SSL))
	if err != nil {
		return nil, fmt.Errorf("error creating a new session : %v", err)
	}

	return &Redshift{
		Client: db,
		Config: *cfg,
	}, nil
}

func (r *Redshift) Schema(table string) (types.Table, error) {

	query := fmt.Sprintf(Redshift_Schema_query, r.Config.DatabaseName, r.Config.Schema, table)
	ctx := context.Background()
	rows, err := r.Client.QueryContext(ctx, query)
	if err != nil {
		return types.Table{}, fmt.Errorf("error executing query: %v", err)
	}

	var columns []types.Column
	for rows.Next() {
		var column types.Column
		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&column.IsNullable,
			&column.DefaultValue,
			&column.CharacterMaximumLength,
			&column.OrdinalPosition,
			&column.Visibility,
			&column.IsPrimary,
			&column.IsUpdatable,
		); err != nil {
			return types.Table{}, fmt.Errorf("error scanning rows: %v", err)
		}
		column.Metatags = []string{}
		column.Metatags = append(column.Metatags, column.Name)
		columns = append(columns, column)

	}

	if err := rows.Err(); err != nil {
		return types.Table{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    []string{},
	}, nil
}

func (r *Redshift) Tables(databaseName string) ([]string, error) {
	ctx := context.Background()
	query := fmt.Sprintf(Redshift_Tables_query, databaseName)

	res, err := r.Client.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	var tables []string
	for res.Next() {
		var table string
		if err := res.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning result: %v", err)
		}
		tables = append(tables, table)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over result: %v", err)
	}

	return tables, nil

}

func (r *Redshift) Execute(query string) ([]byte, error) {
	ctx := context.Background()
	rows, err := r.Client.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// getting the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	// Scan the result into a slice of slices
	var results [][]interface{}
	for rows.Next() {
		// create a slice of values and pointers
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			//  create a slice of pointers to the values
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

func (r *Redshift) GenerateCreateTableQuery(table types.Table) string {
	query := fmt.Sprintf("CREATE TABLE %s.%s.%s (", r.Config.DatabaseName, r.Config.Schema, table.Name)
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + colType

		if column.IsPrimary {
			query += " PRIMARY KEY"
			if column.AutoIncrement {
				query += fmt.Sprintf(" IDENTITY(%v, %v)", column.IdentitySeed, column.IdentityStep)
			}
		}

		if column.IsNullable == "NO" {
			query += " NOT NULL"
		}

		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ");"
	return query
}
