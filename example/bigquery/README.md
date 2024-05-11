## Getting Started with Bigquery and Go

This guide will walk you through the process of setting up a Bigquery and running a Go application that interacts with it.


### Step 1: Set Up BigQuery

To get started, ensure you have access to Google Cloud Platform (GCP) and have created a BigQuery dataset.

### Step 2: Configure Your BigQuery Settings

In your Go application, define your BigQuery database configuration. This includes specifying your GCP project ID and providing the path to your JSON key file for authentication.

### Step 3: Install Xray and Import Dependencies

Make sure you have Xray installed and imported in your Go application. You can do this by adding the necessary import statement in your Go code:

    ```
    import (
    "github.com/thesaas-company/xray"
    "github.com/thesaas-company/xray/config"
    "github.com/thesaas-company/xray/types"
    )
    ```

### Step 4: Use Xray to Connect to BigQuery

Utilize Xray to connect to your BigQuery database. Use the provided code sample in [Integration.md](https://github.com/thesaas-company/xray/tree/main/example/bigquery/integration.md) to see how to set up the connection and execute queries.

### Step 5: Run Your Go Application

Once you have configured your BigQuery settings and integrated Xray into your Go application, run the following commands to ensure your project dependencies are up to date and to execute your application:


    ```
    go mod tidy
    go run main.go
    ```

That's it! You should now be able to inspect and execute queries in your BigQuery database using Xray in your Go application.