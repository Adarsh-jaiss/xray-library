## Integrating Xray with Microsoft SQL Server

This guide illustrates how to use the Xray library to inspect and execute queries in a Microsoft SQL Server (MSSQL) database within a Go application.

### Step-by-Step Guide

1. **Define MSSQL Configuration**

   Add your Microsoft SQL Server configuration to your Go application:

    **Note : Set env variable DB_Password for adding password and pass your password as : `export DB_PASSWORD=your_password`**

    ```go
    cfg := &config.Config{
        Host:       "localhost",
        Username:     "sa",
        Port:         "14330",
        DatabaseName: "master",
    }
    ```

2. **Connect to MSSQL Database**

    Create a new instance of the MSSQL database:

    ```go
    client, err := xray.NewClientWithConfig(&config, types.MSSQL)
    if err != nil {
        panic(err)
    }
    ```

3. **Execute Queries**

    Execute queries against your MSSQL database:

    ```go
    query := "SELECT * FROM my_table" // Specify your SQL query
    result, err := client.Execute(query)
    if err != nil {
        fmt.Printf("Error executing query: %v\n", err)
        return
    }
    fmt.Printf("Query result: %s\n", string(result))
    ```

4. **Retrieve Tables and Schema**

    Retrieve a list of tables in the database and print their schemas:

    ```go
    data, err := client.Tables(cfg.DatabaseName)
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
    ```

5. **Generate Create Table Query**

    Generate and print SQL CREATE TABLE queries for each table in the response slice:

    ```go
    // Iterate over each table in the response slice.
    for _, v := range response {
        // Generate a CREATE TABLE query for the current table.
        query := client.GenerateCreateTableQuery(v)
        // Print the generated query.
        fmt.Println(query)
    }
    ```

### Running the Application

After configuring your MSSQL settings and integrating Xray into your Go application, run the following commands to ensure your project dependencies are up to date and to execute your application:

```sh
go mod tidy
go run main.go
```

That's it! You should now be able to inspect and execute queries in your MSSQL database using Xray in your Go application.

