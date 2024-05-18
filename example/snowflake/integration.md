## Using Xray with Snowflake Integration

This guide demonstrates how to use the Xray library to inspect and execute queries in a Snowflake database in a Go application.

### Step-by-Step Guide

1. **Define Snowflake Configuration**

   Add your Snowflake database configuration to your Go application:

    **Note : Set env variable DB_Password for adding password and pass your password as : `export DB_PASSWORD=your_password`**

   ```go
   config := &config.Config{
       Account:      "account",         // Replace with your Snowflake account name
       Username:     "Username",        // Replace with your Snowflake username
       DatabaseName: "DatabaseName",    // Replace with your Snowflake database name
       Port:         "443",             // Snowflake port (default is 443)
       Warehouse:    "Warehouse_name",  // Replace with your Snowflake warehouse name
       Schema:       "Schema",          // Optional: Replace with your Snowflake schema name if applicable
   }

2. **Connect to Snowflake Database**

    Create a new Xray client and connect to the Snowflake database:

        ```go
        client, err := xray.NewClientWithConfig(config, types.Snowflake)
        if err != nil {
            panic(err)
        }
        fmt.Println("Connected to database")
    

3. **Retrieve Tables and Schema**

    Retrieve a list of tables in the database and print their schemas:

        ```go
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
        ```
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

    Execute queries against your snowflake database:

        ```go
        query := "SELECT * FROM database_name.schema_name.my_table"   // Specify your SQL query
        result, err := client.Execute(query)
        if err != nil {
            fmt.Printf("Error executing query: %v\n", err)
            return
        }
        fmt.Printf("Query result: %s\n", string(result))
 



#### Environment Variables

You can also set Snowflake configurations using environment variables. For example:

    ```
    export SNOWFLAKE_ACCOUNT=account
    export SNOWFLAKE_USERNAME=username
    export SNOWFLAKE_PASSWORD=password
    export SNOWFLAKE_DATABASE=database
    export SNOWFLAKE_WAREHOUSE=warehouse
    ```

#### Running the Application

After configuring your Snowflake settings, run the following commands to execute the application:

    ```
    go mod tidy
    go run main.go

    ```

That's it! You should now be able to inspect and execute queries in your Snowflake database using Xray in your Go application.


This integration.md file provides detailed steps on how to connect to a Snowflake database, retrieve tables and schema, and generate SQL queries in a Go application using Xray.
