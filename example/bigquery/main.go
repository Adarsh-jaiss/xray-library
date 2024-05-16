package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// export GOOGLE_APPLICATION_CREDENTIALS=path/to/secret.json
func main() {
	config := &config.Config{
		ProjectID:    "ProjectID",
		DatabaseName: "DatabaseName",
	}

	client, err := xray.NewClientWithConfig(config, types.BigQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	table := types.Table{
		Dataset: config.DatabaseName,
		Name:    "table",
		Columns: []types.Column{
			{Name: "id", Type: "INT64", IsPrimary: true},
			{Name: "name", Type: "STRING"},
			{Name: "created_at", Type: "TIMESTAMP"},
		},
	}

	query := client.GenerateCreateTableQuery(table)
	fmt.Println(query)
	res, err := client.Execute(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

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
