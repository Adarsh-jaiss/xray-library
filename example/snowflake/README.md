## Running a Snowflake Database with Xray

1. To get your Snowflake database up and running with the Xray library, follow these steps:

2. Create a Snowflake Account: If you haven't already, create a Snowflake account to get started.

3. Configure Snowflake Settings: Open your `main.go` file and configure your Snowflake settings as described in detail in the [Integration.md](https://github.com/thesaas-company/xray/tree/main/example/snowflake/integration.md) file.

4. Set Snowflake Password: In your terminal, set your Snowflake password as an environment variable. Replace `<DB_PASSWORD>` with your actual Snowflake password:


    ```bash
    export DB_PASSWORD=<DB_PASSWORD>
    ```

5. Update Dependencies: Ensure your Go project's dependencies are up to date by running the following command in your terminal:
   
    ```
    go mod tidy
    ```

6. Run Your Application: Finally, execute your `main.go` file with the following command:

    ```
    go run main.go

