# MSSQL Integration with Xray

This repository provides an example of how to use the Xray library to interact with a Microsoft SQL Server (MSSQL) database in a Go application. The example demonstrates how to connect to the database, execute queries, retrieve table schemas, and generate SQL `CREATE TABLE` statements.

## Prerequisites

Before you begin, ensure you have the following:

- Go installed on your machine.
- Access to a Microsoft SQL Server instance.
- You can run it locally using this command :

    ```
    docker run -p 14330:1433 -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=root' -d chriseaton/adventureworks:latest
    ```
- Environment variable `DB_PASSWORD` set to your database password.

## Getting Started

### Create a Go Application

Start by creating a main.go file in your Go project and add the provided [example](https://github.com/thesaas-company/xray/tree/main/example/mssql.main.go) in `main.go` code into it and checkout [Integration.md](https://github.com/thesaas-company/xray/tree/main/example/mssql/integration.md) for more info.

### Run the Application

Once you have added the example code to your main.go file, execute the following commands in your terminal:

```
go mod tidy
go run main.go
```
That's it! You should now be able to inspect and execute queries in your MSSQL database using Xray in your Go application.