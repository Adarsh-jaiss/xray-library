## Using Xray with MySQL Integration

This guide demonstrates how to use the Xray library to inspect a MySQL database in a Go application.

### Introduction to Xray

Xray is a library that provides tools for inspecting and analyzing various types of databases. In this example, we'll use Xray to connect to a MySQL database, retrieve the schema for each table, and generate SQL queries.

### Step-by-Step Guide

1. **Define Database Configuration**

   ```go
   // Define your database configuration here
   dbConfig := &config.Config{
       Host:         "127.0.0.1",
       DatabaseName: "employees",
       Username:     "root",
       Port:         "3306",
       SSL:          "false",
   }

**Note** : Ensure to replace the placeholders with your actual database configuration.

2. **Connect to MySQL Database**
   
    ```go
    client, err := xray.NewClientWithConfig(dbConfig, types.MySQL)
    if err != nil {
        panic(err)
    }
3. **Retrieve Table Schema**

    ```go
    data, err := client.Tables(dbConfig.DatabaseName)
    if err != nil {
        panic(err)
    }
4. **Print Table Schema**

    ```go
    for _, tableName := range data {
        table, err := client.Schema(tableName)
        if err != nil {
            panic(err)
        }
        fmt.Println(table)
    }

5. **Define and Generate SQL Query for a New Table**

    ```go
    newTable := types.Table{
    Name: "user",
    Columns: []types.Column{
            {
                Name:       "id",
                Type:       "int",
                IsNullable: "NO",
                IsPrimary:  true,
            },
            {
                Name:       "name",
                Type:       "varchar(255)",
                IsNullable: "NO",
                IsUnique:   true,
            },
            {
                Name:       "age",
                Type:       "int",
                IsNullable: "YES",
            },
        },
    }

    query := client.GenerateCreateTableQuery(newTable)
    fmt.Println(query)


### Environment Variables

- You can also configure the MySQL connection using environment variables. Set the following environment variables before running your application:

```
export DB_HOST=127.0.0.1
export DB_NAME=employees
export DB_USERNAME=root
export DB_PORT=3306
export DB_SSL=false
```

### Running the Application

After setting up the configuration, run the following commands to execute the application:

```
go mod tidy
go run main.go
```
