package main

import (
	"fmt"

	_ "github.com/snowflakedb/gosnowflake"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// export DB Passowrd, Export root=DB_PASSWORD
func main() {
	config := &config.Config{
		Account:      "account",
		Username:     "Username",
		DatabaseName: "DatabaseName",
		Port:         "443",
		Warehouse:    "Wareshousw_name",
		Schema:       "Schema", // optional
	}

	client, err := xray.NewClientWithConfig(config, types.Snowflake)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	data, err := client.Tables(config.DatabaseName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables :", data)

	var response []types.Table
	for _, v := range data {
		table, err := client.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)

	for _, v := range response {
		query := client.GenerateCreateTableQuery(v)
		fmt.Println(query)
	}
}
