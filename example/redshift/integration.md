## Integrating Xray with Amazon Redshift

This guide illustrates how to use the Xray library to inspect and execute queries in an Amazon Redshift database within a Go application.

### Step-by-Step Guide
1. **Define Redshift Configuration**

    Start by defining the configuration for your Amazon Redshift instance in your Go application:

    ```go
    cfg := &config.Config{
        AWS: config.AWS{
            Region:          "us-west-2",    // Specify your AWS region
            AccessKey:       "access_key",   // Add your AWS access key
            SecretAccessKey: "secret_access",// Add your AWS secret access key
        },
        DatabaseName: "my-database",   // Specify your Redshift database name
        Schema:       "my-schema",     // Specify your Redshift schema name
    }
    ```

2. **Connect to Redshift Database**

    Create a new instance of the Amazon Redshift database:

    ```go
    rs, err := redshift.NewRedshiftWithConfig(cfg)
    if err != nil {
        fmt.Printf("Error creating Redshift instance: %v\n", err)
        return
    }
    ```

3. **Execute Queries**

    Execute queries against your Redshift database:

        ```go
        query := "SELECT * FROM my_table"   // Specify your SQL query
        result, err := rs.Execute(query)
        if err != nil {
            fmt.Printf("Error executing query: %v\n", err)
            return
        }
        fmt.Printf("Query result: %s\n", string(result))
        ```

5. **Retrieve Tables and Schema**

    Retrieve a list of tables in the database and print their schemas:

    ```go
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
    ```

#### Running the Application

After configuring your Amazon Redshift settings and integrating Xray into your Go application, run the following commands to ensure your project dependencies are up to date and to execute your application:

    go mod tidy
    go run main.go
    
That's it! You should now be able to inspect and execute queries in your Amazon Redshift database using Xray in your Go application.
