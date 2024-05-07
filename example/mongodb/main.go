package mongodb

import (
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/mongoDB"
	"github.com/thesaas-company/xray/types"
)

// export DB_PASSWORD=your_password
func main() {
	config := &config.Config{
		Username:     "admin",
		Host:         "localhost",
		Port:         "27017",
		DatabaseName: "my-database",
	}

	// Create a new MongoDB client
	client, err := mongodb.NewMongoDBWithConfig(config)
	if err != nil {
		panic(fmt.Errorf("error creating new MongoDB client: %v", err))
	}

	// Get the collections in the database
	collections, err := client.Tables(config.DatabaseName)
	if err != nil {
		panic(fmt.Errorf("error getting collections: %v", err))
	}
	fmt.Println(collections)

	// Get the schema of a collection
	var response []types.Table
	for _, v := range collections {
		table, err := client.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)
}

