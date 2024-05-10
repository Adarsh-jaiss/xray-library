package main

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

func main() {
	config := &config.Config{
		ProjectID:    "ProjectID",
		JSONKeyPath:  "/JSONKeyPath",
		DatabaseName: "Database_Name",
	}

	client, err := xray.NewClientWithConfig(config, types.BigQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

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

	query := client.GenerateCreateTableQuery(table)
	fmt.Println(query)

	tables, err := client.Tables(config.DatabaseName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables :", tables)

	var response []types.Table
	for _, v := range tables {
		table, err := client.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)
}
