package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

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
}
