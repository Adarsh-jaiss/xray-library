## Integrating Xray with Amazon Redshift

This guide illustrates how to use the Xray library to inspect and execute queries in an Amazon Redshift database within a Go application.

### Step-by-Step Guide
1. **Define Redshift Configuration**

    Start by defining the configuration for your Amazon Redshift instance in your Go application:
    
    **Note : Set env variable DB_Password for adding password and pass your password as : `export DB_PASSWORD=your_password`**

    ```go
    cfg := &config.Config{
        Host:         "default-workgroup.587687374938.ap-north-2.redshift-serverless.amazonaws.com"
        Username:     "admin",
		Port:         "5439",
        SSL:          "require",    // 
        DatabaseName: "my-database",   // Specify your Redshift database name
    }
    ```

2. **Connect to Redshift Database**

    Create a new instance of the Amazon Redshift database:

    ```go
    client, err := xray.NewClientWithConfig(cfg,types.Redshift)
	if err != nil {
		fmt.Printf("Error creating Redshift instance: %v\n", err)
		return
	}
    ```

     ```

3. **Retrieve Tables and Schema**

    Retrieve a list of tables in the database and print their schemas:

    ```go
    tables, err := client.Tables(cfg.DatabaseName)
    if err != nil {
        fmt.Printf("Error getting tables: %v\n", err)
        return
    }
    fmt.Printf("Tables in database %s: %v\n", cfg.DatabaseName, tables)
    var response []types.Table
    for _, v := range tables {
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

    Execute queries against your Redshift database:

        ```go
        query := "SELECT * FROM my_table"   // Specify your SQL query
        result, err := client.Execute(query)
        if err != nil {
            fmt.Printf("Error executing query: %v\n", err)
            return
        }
        fmt.Printf("Query result: %s\n", string(result))
   

#### Running the Application

After configuring your Amazon Redshift settings and integrating Xray into your Go application, run the following commands to ensure your project dependencies are up to date and to execute your application:

    ```
    go mod tidy
    go run main.go
    
That's it! You should now be able to inspect and execute queries in your Amazon Redshift database using Xray in your Go application.
