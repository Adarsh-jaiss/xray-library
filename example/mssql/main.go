package main

import (
	// "database/sql"
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// export DB_PASSWORD=your_password
func main() {
	config := config.Config{
		Server:       "localhost",
		Username:     "sa",
		Port:         "14330",
		DatabaseName: "master",
	}

	client, err := xray.NewClientWithConfig(&config, types.MSSQL)
	if err != nil {
		panic(err)
	}

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
