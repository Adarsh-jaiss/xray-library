## Using Xray with Postgres Integration

This guide demonstrates how to use the Xray library to inspect a Postgres database in a Go application.

### Introduction to Xray

Xray is a library that provides tools for inspecting and analyzing various types of databases. In this example, we'll use Xray to connect to a Postgres database, retrieve the schema for each table, and print it.

### Step-by-Step Guide

1. **Define Database Configuration**

   ```go
   // Define your database configuration here
   dbConfig := &config.Config{
       Host:         "127.0.0.1",
       DatabaseName: "employees",
       Username:     "root",
       Port:         "5432",
       SSL:          "disable",
   }

***Note** : Ensure to replace the placeholders with your actual database configuration.

2. **Connect to Postgres Database**

    ```go
    client, err := xray.NewClientWithConfig(dbConfig, types.Postgres)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected to database")

3. **Retrieve Table Names**

    ```go
    data, err := client.Tables(dbConfig.DatabaseName)
    if err != nil {
        panic(err)
    }
    fmt.Println("Tables :", data)

4. **Print Table Schema**

    ```go
    var response []types.Table
    for _, tableName := range data {
        table, err := client.Schema(tableName)
        if err != nil {
            panic(err)
        }
        response = append(response, table)
    }
    fmt.Println(response)

#### Environment Variables

- You can also configure the Postgres connection using environment variables. Set the following environment variables before running your application:
    ```
    export DB_HOST=127.0.0.1
    export DB_NAME=employees
    export DB_USERNAME=root
    export DB_PORT=5432
    export DB_SSL=disable

#### Running the Application

After setting up the configuration, run the following commands to execute the application:

    ```
    go mod tidy
    go run main.go
