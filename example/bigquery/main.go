package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

func main() {
	config := &config.Config{
		ProjectID:    "ProjectID",
		JSONKeyPath:  "JSONKeyPath",
		DatabaseName: "Database_Name",
	}

	client, err := xray.NewClientWithConfig(config, types.BigQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

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