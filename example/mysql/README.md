
## Getting Started with MySQL and Go

This guide will walk you through the process of setting up a MySQL server using Docker and running a Go application that interacts with it.

### Step 1: Run a MySQL Server

First, we need to set up a MySQL server. We will use Docker to run a MySQL server with preloaded employees data (for demo). Run the following command in your terminal:

```bash
docker run -d --name mysql-employees \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=college \
  genschsa/mysql-employees
```
**NOTE** :Make sure to replace college with your desired password.

#### Set Password in env variable
Set the password in an environment variable for easier configuration:
```bash
export DB_PASSWORD=college
```

### Step 2: Configure and Run the Go Application

- Define your MySQL database configuration in your Go application.

- Ensure you have Xray installed and imported in your Go application.

- Checkout [Integration.md](https://github.com/thesaas-company/xray/tree/main/example/mysql/integration.md) for a code sample demonstrating how to connect to MySQL using Xray.

- Once you have the code, run the following commands:
```
go mod tidy
go run main.go
```
