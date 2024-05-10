package main

import (
	"database/sql"
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

// export DB_PASSWORD=your_password
func main() {
	config := &config.Config{
		Host:         "127.0.0.1",
		DatabaseName: "employees",
		Username:     "root",
		Port:         "3306",
		SSL:          "false",
	}
	client, err := xray.NewClientWithConfig(config, types.MySQL)
	if err != nil {
		panic(err)
	}
	data, err := client.Tables(config.DatabaseName)
	if err != nil {
		panic(err)
	}
	var response []types.Table
	for _, v := range data {
		table, err := client.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)

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
}
