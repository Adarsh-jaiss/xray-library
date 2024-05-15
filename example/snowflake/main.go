package main

import (
	"database/sql"
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
}
