package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)
// export DB_PASSWORD=your_password
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
	client, err := xray.NewClientWithConfig(cfg, types.Redshift)
	if err != nil {
		fmt.Printf("Error creating Redshift instance: %v\n", err)
		return
	}
	fmt.Println("Connected to database")

	tables, err := client.Tables(cfg.DatabaseName)
	if err != nil {
		fmt.Printf("Error getting tables: %v\n", err)
		return
	}
	fmt.Printf("Tables in database %s: %v\n", cfg.DatabaseName, tables)

	var response []types.Table
	for _, v := range tables {
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
