# Interact with Xray using Xray-CLI

To execute commands in SQL databases directly from shell you can use different configs for different databases and warehouses.
The CLI supports Mysql, Postgres, Microsoft SQL Server, Snowflake, Bigquery and Redhsift.

Simply run this command : `xray` to see all the commands in the CLI application.
If you need more info about a specific command, you can run this  : `xray <COMMAND> --help`.
For example: Currently we only support shell command, which is used for interacting with different SQL databases. you can see the available options or more info via this command : `xray shell --help` or `xray shell -h`. you can use this command to see the available flags and their use case.

To execute query in any database, use this command : 

```
xray shell -t <DATABASE TYPE> -c <Config.yaml file location>
```

if you want the result in verbose mode, you can add the -v flag in the command as well

```
xray shell -t <DATABASE TYPE> -c <Config.yaml file location> -v
```


### Mysql

To run mysql and interact with it, simply run this command :

<!-- Export DB_PASSWORD=YOURPASSWORD -->

```
xray shell -t mysql -c example/mysql/config.yaml
```

For verbose mode : 
```
xray shell -t mysql -c example/mysql/config.yaml -v
```

### Postgres

To run postgres and interact with it, simply run this command :

<!-- Export DB_PASSWORD=YOURPASSWORD -->

```
xray shell -t postgres -c example/postgres/config.yaml
```

For verbose mode : 
```
xray shell -t postgres -c example/postgres/config.yaml -v
```

### MSSQL

To run mssql and interact with it, simply run this command :

<!-- Export DB_PASSWORD=YOURPASSWORD -->

```
xray shell -t mssql -c example/mssql/config.yaml
```

For verbose mode : 
```
xray shell -t mssql -c example/mssql/config.yaml -v
```

### Snowflake

To run snowflake and interact with it, simply run this command :

<!-- Export DB_PASSWORD=YOURPASSWORD -->

```
xray shell -t snowflake -c example/snowflake/config.yaml
```

For verbose mode : 
```
xray shell -t snowflake -c example/snowflake/config.yaml -v
```

### Redhshift

To run redshift and interact with it, simply run this command :

<!-- export DB_PASSWORD=YOURPASSWORD -->

```
xray shell -t redshift -c example/redshift/config.yaml
```

For verbose mode : 
```
xray shell -t redshift -c example/redshift/config.yaml -v
```

### Bigquery

To run bigquery and interact with it, simply run this command :

<!-- export GOOGLE_APPLICATION_CREDENTIALS=path/to/secret.json-->

```
xray shell -t bigquery -c example/bigquery/config.yaml
```

For verbose mode : 
```
xray shell -t bigquery -c example/bigquery/config.yaml -v
```