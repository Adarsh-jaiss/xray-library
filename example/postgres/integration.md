## Using Xray with Postgres Integration

This guide demonstrates how to use the Xray library to inspect a Postgres database in a Go application.

### Introduction to Xray

Xray is a library that provides tools for inspecting and analyzing various types of databases. In this example, we'll use Xray to connect to a Postgres database, retrieve the schema for each table, and print it.

### Step-by-Step Guide

1. **Define Database Configuration**

     Add your Postgres database configuration to your Go application:

    **Note : Set env variable DB_Password for adding password and pass your password as : `export DB_PASSWORD=your_password`**
 
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

1. **Connect to Postgres Database**

    ```go
    client, err := xray.NewClientWithConfig(dbConfig, types.Postgres)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected to database")

2. **Retrieve Table Names**

    ```go
    data, err := client.Tables(dbConfig.DatabaseName)
    if err != nil {
        panic(err)
    }
    fmt.Println("Tables :", data)

3. **Print Table Schema**

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

4. **Generate Create Table Query**
   
   generate and print SQL CREATE TABLE queries for each table in the response slice.

   ```go
    // Iterate over each table in the response slice.
    for _, v := range response {
        // Generate a CREATE TABLE query for the current table.
        query := client.GenerateCreateTableQuery(v)
        // Print the generated query.
        fmt.Println(query)
    }
    ```

5. **Execute Queries**

    Execute queries against your  database:

    ```go
    query := `SELECT * FROM "my_table"` // Specify your SQL query
    result, err := client.Execute(query)
    if err != nil {
        fmt.Printf("Error executing query: %v\n", err)
        return
    }
    fmt.Printf("Query result: %s\n", string(result))
  



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
