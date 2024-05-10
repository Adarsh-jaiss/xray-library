# XRay

XRay is a open source library for database schema extraction and query execution.

## Install 

```bash
go get github.com/thesaas-company/xray@latest
```

## Docs

- Official documentation : [Go Docs](https://pkg.go.dev/github.com/thesaas-company/xray)
- Examples of different database integrations : [Example](./example)

## Getting started 

### Run a MySQL Server

```bash
docker run -d --name mysql-employees \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=college \
  genschsa/mysql-employees
```
- Set Password in env variable
```bash
export DB_PASSWORD=college
```
- Use Xray to inspect your MySQL database, by creating a main.go file and adding the below example into it.
  
```
package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)


func main() {
  // Define your configs here
	config := &config.Config{
		Host:         "127.0.0.1", // <YOUR HOST IP>
		DatabaseName: "employees", // <YOUF DATABASE NAME>
		Username:     "root",      // <MYSQL USERNAME>
		Port:         "3306",      // <MYSQL PORT>
		SSL:          "false",     // <SSL config>
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

```
and then run these commands : 
```
go mod tidy
go run main.go
```

### Run a Postgres Database

- To use the postgres docker instance you can use this below command with your configs

```
docker run -d --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=employees \
  -v "$(pwd)/data:/var/lib/postgresql/data" \
  postgres:13.2-alpine
```

- To inspect and execute queries in Postgres database, Create a main.go file and adding the below example into it.

```
package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

func main() {
  // Define your configs here
	config := &config.Config{
		Host:         "127.0.0.1",  // <YOUR POSTGRES HOST IP>
		DatabaseName: "employees",  // <YOUR POSTGRES DATABASE NAME>
		Username:     "root",       // <YOUR POSTGRES USERNAME>
		Port:         "5432",       // <POSTGRES PORT NO>
		SSL:          "disable",    // <SSL (enable/disable)>
	}
	client, err := xray.NewClientWithConfig(config, types.Postgres)
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

}

```
- and then run these commands : 

```
go mod tidy
go run main.go
```

### Run a Snowflake Database

- To inspect and execute queries in snowflake database, Create a main.go file and adding the below example into it.


```
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
    // Add your snowflake configs below
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
}
```

- and then run these commands: 

```
go mod tidy
go run main.go
```

#### Run Bigquery 

- To inspect and execute queries in bigquery database, Create a main.go file and adding the below example into it.

```
package main

import (
	"fmt"

	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

func main() {
  // Add your bigquery configs here
	config := &config.Config{
		ProjectID:    "ProjectID",
		JSONKeyPath:  "JSONKeyPath", // create a secrets.json file, add your configs and give the file path here
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
```
- and then run these commands: 

```
go mod tidy
go run main.go
```

### Run Redhsift

- To inspect and execute queries in a Redshift database, Create a main.go file and adding the below example into it.

```
package main

import (
	"fmt"

	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/databases/redshift"
	"github.com/thesaas-company/xray/types"
)

func main() {
	// Define the configuration for the Redshift instance
	cfg := &config.Config{
		AWS: config.AWS{
			Region:          "us-west-2",
			AccessKey:       "access_key",
			SecretAccessKey: "secret_access",
		},
		DatabaseName: "my-database",
		Schema:       "my-schema",
	}

	// Create a new Redshift instance
	rs, err := redshift.NewRedshiftWithConfig(cfg)
	if err != nil {
		fmt.Printf("Error creating Redshift instance: %v\n", err)
		return
	}

	query := "SELECT * FROM my_table"
	result, err := rs.Execute(query)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	fmt.Printf("Query result: %s\n", string(result))

	tables, err := rs.Tables(cfg.DatabaseName)
	if err != nil {
		fmt.Printf("Error getting tables: %v\n", err)
		return
	}
	fmt.Printf("Tables in database %s: %v\n", cfg.DatabaseName, tables)
	var response []types.Table
	for _, v := range tables {
		table, err := rs.Schema(v)
		if err != nil {
			panic(err)
		}
		response = append(response, table)
	}
	fmt.Println(response)
}
```

- and then run these commands: 

```
go mod tidy
go run main.go
```


## Maintainer
- [@Adarsh Jaiswal](https://github.com/adarsh-jaiss)
- [@tqindia](https://github.com/tqindia)
