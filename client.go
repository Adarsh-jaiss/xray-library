package xray

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/bigquery"

	"github.com/thesaas-company/xray/databases/mysql"
	"github.com/thesaas-company/xray/databases/postgres"
	"github.com/thesaas-company/xray/databases/redshift"
	"github.com/thesaas-company/xray/databases/snowflake"
	"github.com/thesaas-company/xray/logger"
	"github.com/thesaas-company/xray/types"
)

// NewClientWithConfig creates a new SQL client with the given configuration and database type.
// It returns an error if the database type is not supported or if there is a problem creating the client.
func NewClientWithConfig(dbConfig *config.Config, dbType types.DbType) (types.ISQL, error) {
	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQLWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgresWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Snowflake:
		sqlClient, err := snowflake.NewSnowflakeWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.BigQuery:
		bigqueryClient, err := bigquery.NewBigQueryWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(bigqueryClient), nil
	case types.Redshift:
		redshiftClient, err := redshift.NewRedshiftWithConfig(dbConfig)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(redshiftClient), nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// NewClient creates a new SQL client with the given database client and database type.
// It returns an error if the database type is not supported or if there is a problem creating the client.
func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {

	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQL(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgres(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Snowflake:
		sqlClient, err := snowflake.NewSnowflake(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.BigQuery:
		BigQueryClient, err := bigquery.NewBigQuery(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(BigQueryClient), nil
	case types.Redshift:
		redshiftClient, err := redshift.NewRedshift(dbClient)
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(redshiftClient), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
