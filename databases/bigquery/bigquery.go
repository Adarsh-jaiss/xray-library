package bigquery

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// The BigQuery struct is responsible for holding the BigQuery client and configuration.
type BigQuery struct {
	Client *bigquery.Client
	Config *config.Config
}

func NewBigQuery(client *bigquery.Client) (types.ISQL, error) {
	return &BigQuery{
		Client: client,
		Config: &config.Config{},
	}, nil
}

func NewBigQueryWithConfig(cfg *config.Config) (types.ISQL, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.JSONKeyPath))
	if err != nil {
		return nil, err
	}

	return &BigQuery{
		Client: client,
		Config: cfg,
	}, nil

}

// this function extarcts the schema of a table in BigQuery.
// It takes table name as input and returns a Table struct and an error.
func (b *BigQuery) Schema(table string) (types.Table, error) {
	ctx := context.Background()
	var schema types.Table

	tableRef := b.Client.Dataset(b.Config.DatabaseName).Table(table)
	schemaInfo, err := tableRef.Metadata(ctx)
	if err != nil {
		return types.Table{}, fmt.Errorf("error getting table metadata: %v", err)
	}

	schema.Name = schemaInfo.Name
	schema.Description = schemaInfo.Description
	schema.Columns = make([]types.Column, len(schemaInfo.Schema))
	for i, fieldSchema := range schemaInfo.Schema {
		schema.Columns[i] = types.Column{
			Name:        fieldSchema.Name,
			Type:        string(fieldSchema.Type),
			Description: fieldSchema.Description,
			CharacterMaximumLength: sql.NullInt64{
				Int64: fieldSchema.MaxLength,
				Valid: fieldSchema.MaxLength != 0,
			},
			DefaultValue: sql.NullString{
				String: fieldSchema.DefaultValueExpression,
				Valid:  fieldSchema.DefaultValueExpression != "",
			},
		}
	}

	fmt.Println(schema)
	return schema, nil
}

// Execute runs a query in BigQuery and returns the results as a byte slice.
// It takes a query string as input and returns a byte slice and an error.
func (b *BigQuery) Execute(query string) ([]byte, error) {
	ctx := context.Background()
	q := b.Client.Query(query)

	exe, err := q.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("error running query")
	}

	status, err := exe.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for query to complete: %v", err)
	}

	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("expected nil, found err: %v", err)
	}

	it, err := exe.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading query results: %v", err)
	}

	var result []map[string]interface{}
	for {
		var values map[string]interface{}
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		result = append(result, values)
	}

	// Marshal the rows into JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query results: %v", err)
	}

	return jsonData, nil
}

// Tables returns a list of tables in a dataset.
// It takes a dataset name as input and returns a slice of strings and an error.

func (b *BigQuery) Tables(Dataset string) ([]string, error) {
	res := b.Client.Dataset(Dataset).Tables(context.Background())
	var tables []string

	for {
		table, err := res.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error scanning dataset")
		}
		tables = append(tables, table.TableID)
	}

	return tables, nil
}
