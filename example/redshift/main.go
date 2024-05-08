package main

import (
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/redshift"
	"github.com/thesaas-company/xray/types"
)

func main() {
	// Define the configuration for the Redshift instance
	cfg := &config.Config{
		AWS: config.AWS{
			Region:          "us-west-2",
			AccessKey:       "access_key",
			SecretAccessKey: "secret_access",
		},
		DatabaseName: "my-database",
		Schema:       "my-schema",
	}

	// Create a new Redshift instance
	rs, err := redshift.NewRedshiftWithConfig(cfg)
	if err != nil {
		fmt.Printf("Error creating Redshift instance: %v\n", err)
		return
	}

	query := "SELECT * FROM my_table"
	result, err := rs.Execute(query)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	fmt.Printf("Query result: %s\n", string(result))

	tables, err := rs.Tables(cfg.DatabaseName)
	if err != nil {
		fmt.Printf("Error getting tables: %v\n", err)
		return
	}
	fmt.Printf("Tables in database %s: %v\n", cfg.DatabaseName, tables)
	var response []types.Table
	for _, v := range tables {
		table, err := rs.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)
}
