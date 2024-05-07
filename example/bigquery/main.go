package bigquery

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

func main() {
	config := &config.Config{
		ProjectID:   "project-id",
		JSONKeyPath: "path/to/key.json",
	}

	client, err := xray.NewClientWithConfig(config, types.BigQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	tables, err := client.Tables("dataset")
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
