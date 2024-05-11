package redshift

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// Redshift is a struct that represents a Redshift database.
type Redshift struct {
	Client *redshift.Redshift
	config *config.Config
}

// Redshift_List_Tables_query is the SQL query used to list all tables in a schema in Redshift.
const (
	Redshift_List_Tables_query = "SELECT *FROM svv_all_tables WHERE database_name = 'tickit_db' = '%s';"
	Redshift_Schema_query      = "SHOW COLUMNS FROM TABLE %s.%s.%s;"
)

// NewRedshift creates a new Redshift instance.
func NewRedshift(client *redshift.Redshift) (types.ISQL, error) {
	return &Redshift{
		Client: client,
		config: &config.Config{},
	}, nil
}

// NewRedshiftWithConfig creates a new Redshift instance with the provided configuration.
func NewRedshiftWithConfig(cfg *config.Config) (types.ISQL, error) {
	// Create a new AWS session
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWS.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AWS.AccessKey, cfg.AWS.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new session: %v", err)
	}

	// Create a new Redshift client
	client := redshift.New(session)

	// Return a new Redshift instance
	return &Redshift{
		Client: client,
		config: cfg,
	}, nil

}

func (r *Redshift) RedshiftAPIService(config *config.Config, query string) (*redshiftdataapiservice.RedshiftDataAPIService, *redshiftdataapiservice.ExecuteStatementInput) {
	// config.AWS.SecretArn = "arn:aws:secretsmanager:us-west-2:123456789012:secret:My	DBSecret-a1b2c3"
	// config.AWS.ClusterIdentifier = "my-cluster"
	// config.DatabaseName = "my-database"
	svc := redshiftdataapiservice.New(session.Must(session.NewSession(&aws.Config{})))
	input := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(r.config.AWS.ClusterIdentifier),
		Database:          aws.String(r.config.DatabaseName),
		SecretArn:         aws.String(r.config.AWS.SecretArn),
		Sql:               aws.String(query),
	}

	return svc, input
}

// Schema returns the schema of a table in Redshift.
// It takes a table name as input and returns a Table struct and an error.
func (r *Redshift) Schema(table string) (types.Table, error) {
	query := fmt.Sprintf(Redshift_Schema_query, r.config.DatabaseName, r.config.Schema, table)
	svc, input := r.RedshiftAPIService(r.config, query)
	result, err := svc.ExecuteStatement(input)
	if err != nil {
		return types.Table{}, fmt.Errorf("error executing statement: %v", err)
	}

	getStatementResultInput := &redshiftdataapiservice.GetStatementResultInput{
		Id: result.Id,
	}

	getStatementResultOutput, err := svc.GetStatementResult(getStatementResultInput)
	if err != nil {
		return types.Table{}, fmt.Errorf("error getting statement result: %v", err)
	}

	var columns []types.Column
	for _, record := range getStatementResultOutput.Records {
		for _, field := range record {
			if field.StringValue != nil {
				columns = append(columns, types.Column{Name: *field.StringValue})

			}
		}
	}

	return types.Table{
		Name:        table,
		Columns:     columns,
		ColumnCount: int64(len(columns)),
		Description: "",
		Metatags:    nil,
	}, nil

}

// Execute executes a query on Redshift.
// It takes a query string as input and returns the result as a JSON byte slice and an error.
func (r *Redshift) Execute(query string) ([]byte, error) {
	svc, input := r.RedshiftAPIService(r.config, query) // create a new Redshift API service
	result, err := svc.ExecuteStatement(input)          // execute the statement
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	// Get the result
	getStatementResultInput := &redshiftdataapiservice.GetStatementResultInput{
		Id: result.Id,
	}

	getStatementResultOutput, err := svc.GetStatementResult(getStatementResultInput)
	if err != nil {
		return nil, fmt.Errorf("error getting statement result: %v", err)
	}

	// Process the result and convert it to JSON
	jsonData, err := json.Marshal(getStatementResultOutput)
	if err != nil {
		return nil, fmt.Errorf("error converting result to JSON: %v", err)
	}

	return jsonData, nil
}

// Tables returns a list of tables in the specified schema.
// It takes a schema name as input and returns a slice of strings and an error.
func (r *Redshift) Tables(DatabaseName string) ([]string, error) {
	query := fmt.Sprintf(Redshift_List_Tables_query, DatabaseName)
	svc, input := r.RedshiftAPIService(r.config, query)

	res, err := svc.ExecuteStatement(input)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	// Get the result
	getStatementResultInput := &redshiftdataapiservice.GetStatementResultInput{
		Id: res.Id,
	}

	getStatementResultOutput, err := svc.GetStatementResult(getStatementResultInput)
	if err != nil {
		return nil, fmt.Errorf("error getting statement result: %v", err)
	}

	var tables []string
	for _, record := range getStatementResultOutput.Records {
		for _, field := range record {
			if field.StringValue != nil {
				tables = append(tables, *field.StringValue)
			}
		}
	}
	return tables, nil
}

func (r *Redshift) GenerateCreateTableQuery(table types.Table) string {
	query := fmt.Sprintf("CREATE TABLE %s.%s.%s (", r.config.DatabaseName, r.config.Schema, table.Name)
	for i, column := range table.Columns {
		colType := strings.ToUpper(column.Type)
		query += column.Name + " " + colType
		if column.IsPrimary {
			query += " PRIMARY KEY"
		}
		if column.AutoIncrement {
			query += " IDENTITY(1,1)"
		}

		if column.DefaultValue.Valid {
			query += fmt.Sprintf(" DEFAULT %s", column.DefaultValue.String)
		}
		if column.IsUnique.String == "YES" {
			query += " UNIQUE"
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
