package config

// Add Logging, You can use any lib - DONE!!!

// Once we are done with mysql and postgres, Let's rethink about the config structure

// Config holds the configuration details for various databases.
type Config struct {
	// Host is the database host URL.
	Host string `yaml:"host" pflag:",Database host url"`

	// Username is the database username.
	Username string `yaml:"username" pflag:",Database username"`

	// DatabaseName is the name of the database.
	DatabaseName string `yaml:"database" pflag:",Database name"`

	// Port is the database port.
	Port string `yaml:"port" pflag:",Database Port"`

	// SSL is used to enable or disable SSL for the database connection.
	SSL string `yaml:"ssl" pflag:",Database ssl enable/disable"`

	// ProjectID is the BigQuery project ID.
	ProjectID string `yaml:"project_id" pflag:",BigQuery project ID"`

	// Warehouse is the Snowflake warehouse.
	Warehouse string `yaml:"warehouse" pflag:",Snowflake warehouse"`

	// Schema is the Snowflake database schema.
	Schema string `yaml:"schema" pflag:",Snowflake database schema"`

	// Account is the Snowflake account ID.
	Account string `yaml:"account" pflag:",Snowflake account ID"`

	// Debug is used to enable or disable debug mode.
	Debug bool `yaml:"debug" pflag:",Debug mode"`

	// Region is the AWS region.
	Region string `yaml:"region" pflag:",AWS region"`

	// AccountID is the AWS account ID.
	AccountID string `yaml:"account_id" pflag:",AWS account ID"`

	// SecretName is the AWS secret name.
	SecretName string `yaml:"secret_name" pflag:",AWS secret name"`

	// AWS holds the AWS configuration details.
	// AWS AWS `yaml:"aws"`
}

// type AWS struct {
// 	// Region is the AWS region.
// 	Region string `yaml:"region" pflag:",AWS region"`

// 	// AccessKey is the AWS access key.
// 	AccessKey string `yaml:"access_key" pflag:",AWS access key"`

// 	// SecretKey is the AWS secret key.
// 	SecretAccessKey string `yaml:"secret_key" pflag:",AWS secret key"`

// 	// ClusterIdentifier is the AWS cluster identifier.
// 	ClusterIdentifier string `yaml:"cluster_identifier" pflag:",AWS cluster identifier"`

// 	// SecretArn is the AWS secret ARN.
// 	SecretArn string `yaml:"secret_arn" pflag:",AWS secret ARN"`
// }
