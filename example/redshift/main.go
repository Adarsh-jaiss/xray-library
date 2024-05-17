package main

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/redshift"
	"github.com/thesaas-company/xray/types"
)

func main() {
	// Define the configuration for the Redshift instance
	cfg := &config.Config{
		Host:         "default-workgroup.609973658768.ap-south-1.redshift-serverless.amazonaws.com",
		DatabaseName: "dev",
		Username:     "admin",
		Port:         "5439",
		SSL:          "require",
	}

	// Create a new Redshift instance
	rs, err := redshift.NewRedshiftWithConfig(cfg)
	if err != nil {
		fmt.Printf("Error creating Redshift instance: %v\n", err)
		return
	}

	table := types.Table{
		Name: "user",
		Columns: []types.Column{
			{
				Name:         "id",
				Type:         "int",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    true,
				IsUnique:     sql.NullString{String: "YES", Valid: true},
			},
			{
				Name:         "name",
				Type:         "varchar(255)",
				IsNullable:   "NO",
				DefaultValue: sql.NullString{String: "", Valid: false},
				IsPrimary:    false,
				IsUnique:     sql.NullString{String: "NO", Valid: true},
			},
			{
				Name:       "age",
				Type:       "int",
				IsNullable: "YES",
			},
		},
	}

	res := rs.GenerateCreateTableQuery(table)
	fmt.Println(res)

	query := "SELECT * FROM sales limit 10"
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
