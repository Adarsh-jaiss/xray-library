## Using Xray with BigQuery Integration
This guide demonstrates how to use the Xray library to inspect and execute queries in a BigQuery database in a Go application.

#### Step-by-Step Guide

1. **Define BigQuery Configuration**

    Add your BigQuery database configuration to your Go application:

    **Note : Set env variable for password and pass your password as : export GOOGLE_APPLICATION_CREDENTIALS=path/to/secret.json**

    ```go
    config := &config.Config{
        ProjectID:    "ProjectID",      // Replace with your BigQuery project ID
        DatabaseName: "Database_Name",  // Replace with your BigQuery database name
    }

    Ensure you replace `"ProjectID"`, , and `"Database_Name"` with your actual BigQuery project ID, and database name respectively.

2. **Connect to BigQuery Database**
    
    Create a new Xray client and connect to the BigQuery database:

    ```go
    client, err := xray.NewClientWithConfig(config, types.BigQuery)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected to database")

3. **Retrieve Tables and Schema**

    Retrieve a list of tables in the database and print their schemas:

    ```go
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

4. **Generate Create Table Query**
   
   generate and print SQL CREATE TABLE queries for each table in the response slice.

   ```go
    // Iterate over each table in the response slice.
    for _, v := range response {
        // Generate a CREATE TABLE query for the current table.
        query := rs.GenerateCreateTableQuery(v)
        // Print the generated query.
        fmt.Println(query)
    }
    ```
5. **Execute Queries**

    Execute queries against your  database:

    ```go
    query := "SELECT * FROM `project_id.dataset_id.my_table`" // Specify your SQL query
    result, err := client.Execute(query)
    if err != nil {
        fmt.Printf("Error executing query: %v\n", err)
        return
    }
    fmt.Printf("Query result: %s\n", string(result))
  

#### Running the Application
    After configuring your BigQuery settings, run the following commands to execute the application:

    ```
    go mod tidy
    go run main.go
    ```
    
    That's it! You should now be able to inspect and execute queries in your BigQuery database using Xray in your Go application.

    This integration.md file provides detailed steps on how to connect to a BigQuery database, retrieve tables and schema, and execute queries in a Go application using Xray.
