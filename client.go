package xray

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/bigquery"
	"github.com/thesaas-company/xray/databases/mssql"

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
	// Create a new SQL client based on the database type
	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQLWithConfig(dbConfig) // NewMySQLWithConfig is a SQL client that connects to a MySQL database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgresWithConfig(dbConfig) // NewPostgresWithConfig is a SQL client that connects to a Postgres database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Snowflake:
		sqlClient, err := snowflake.NewSnowflakeWithConfig(dbConfig) // NewSnowflakeWithConfig is a SQL client that connects to a Snowflake database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.BigQuery:
		bigqueryClient, err := bigquery.NewBigQueryWithConfig(dbConfig) // NewBigQueryWithConfig is a SQL client that connects to a BigQuery database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(bigqueryClient), nil
	case types.Redshift:
		redshiftClient, err := redshift.NewRedshiftWithConfig(dbConfig) // NewRedshiftWithConfig is a SQL client that connects to a Redshift database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(redshiftClient), nil
	case types.MSSQL:
		mssqlClient, err := mssql.NewMSSQLFromConfig(dbConfig) // NewMSSQLFromConfig is a SQL client that connects to a MSSQL database using the given configuration.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(mssqlClient), nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType) // Return an error if the database type is not supported.
	}
}

// NewClient creates a new SQL client with the given database client and database type.
// It returns an error if the database type is not supported or if there is a problem creating the client.
func NewClient(dbClient *sql.DB, dbType types.DbType) (types.ISQL, error) {
	// Create a new SQL client based on the database type
	switch dbType {
	case types.MySQL:
		sqlClient, err := mysql.NewMySQL(dbClient) // NewMySQL is a SQL client that connects to a MySQL database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Postgres:
		sqlClient, err := postgres.NewPostgres(dbClient) // NewPostgres is a SQL client that connects to a Postgres database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.Snowflake:
		sqlClient, err := snowflake.NewSnowflake(dbClient) // NewSnowflake is a SQL client that connects to a Snowflake database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(sqlClient), nil
	case types.BigQuery:
		BigQueryClient, err := bigquery.NewBigQuery(dbClient) // NewBigQuery is a SQL client that connects to a BigQuery database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(BigQueryClient), nil
	case types.Redshift:
		redshiftClient, err := redshift.NewRedshift(dbClient) // NewRedshift is a SQL client that connects to a Redshift database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(redshiftClient), nil
	case types.MSSQL:
		mssqlClient, err := mssql.NewMSSQL(dbClient) // NewMSSQL is a SQL client that connects to a MSSQL database using the given database client.
		if err != nil {
			return nil, err
		}
		return logger.NewLogger(mssqlClient), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType) // Return an error if the database type is not supported.
	}
}
