## Setting Up a Postgres Database

This guide will walk you through the process of setting up a Postgres database using Docker and running a Go application that interacts with it.

### Step 1: Run a Postgres Database

To use the Postgres Docker instance, run the following command in your terminal:

```bash
docker run -d --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=employees \
  -v "$(pwd)/data:/var/lib/postgresql/data" \
  postgres:13.2-alpine

```

**Note** : Replace root with your desired password.

#### Set Password in Environment Variable

Set the password in an environment variable for easier configuration:
  
  ```bash
  export DB_PASSWORD=root
  ```
  
### Step 2: Configure and Run the Go Application

- Define your Postgres database configuration in your Go application.

- Ensure you have Xray installed and imported in your Go application.

- Checkout [Integration.md](https://github.com/thesaas-company/xray/tree/main/example/postgres/integration.md) for a detailed doc and [main.go](https://github.com/thesaas-company/xray/tree/main/example/postgres/main.go) for code sample, demonstrating how to connect to Postgres using Xray.

- Once you have the code, run the following commands:
  
  ```
  go mod tidy
  go run main.go
  ```




