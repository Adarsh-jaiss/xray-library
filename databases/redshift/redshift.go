package redshift

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

type Redshift struct {
	Client *redshift.Redshift
}

const (
	Redshift_List_Tables_query = "SELECT *FROM svv_all_tables WHERE database_name = 'tickit_db' = '%s';"
	Redshift_Schema_query      = "SHOW COLUMNS FROM TABLE %s.%s.%s;"
)

// NewRedshift creates a new Redshift instance.
func NewRedshift(client *redshift.Redshift) (types.ISQL, error) {
	return &Redshift{
		Client: client,
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
	}, nil

}

// Execute executes a query on Redshift.
func (r *Redshift) RedshiftAPIService(config *config.Config, query string) (*redshiftdataapiservice.RedshiftDataAPIService, *redshiftdataapiservice.ExecuteStatementInput) {
	config.AWS.SecretArn = "arn:aws:secretsmanager:us-west-2:123456789012:secret:My	DBSecret-a1b2c3"
	config.AWS.ClusterIdentifier = "my-cluster"
	config.DatabaseName = "my-database"
	svc := redshiftdataapiservice.New(session.New())
	input := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(config.AWS.ClusterIdentifier),
		Database:          aws.String(config.DatabaseName),
		SecretArn:         aws.String(config.AWS.SecretArn),
		Sql:               aws.String(query),
	}

	return svc, input
}

// Schema returns the schema of a table in Redshift.
// It takes a table name as input and returns a Table struct and an error.
func (r *Redshift) Schema(table string) (types.Table, error) {
	var redshiftConfig config.Config
	query := fmt.Sprintf(Redshift_Schema_query, redshiftConfig.DatabaseName, redshiftConfig.Schema, table)
	svc, input := r.RedshiftAPIService(&redshiftConfig, query)
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

func (r *Redshift) Execute(query string) ([]byte, error) {
	var config config.Config
	svc, input := r.RedshiftAPIService(&config, query)
	result, err := svc.ExecuteStatement(input)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

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
	var config config.Config
	query := fmt.Sprintf(Redshift_List_Tables_query, DatabaseName)
	svc, input := r.RedshiftAPIService(&config, query)

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
