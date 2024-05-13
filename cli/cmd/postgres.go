/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
	"os"
)

var pgClient types.ISQL

// postgresCmd represents the postgres command
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Interact with PostgreSQL databases",
	Long:  `This command allows you to interact with PostgreSQL databases. You can use this command to connect to a PostgreSQL database and run queries.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		user, _ := cmd.Flags().GetString("user")
		database, _ := cmd.Flags().GetString("database")
		ssl, _ := cmd.Flags().GetString("ssl")

		config := &config.Config{
			Host:         host,
			Port:         port,
			Username:     user,
			DatabaseName: database,
			SSL:          ssl,
		}

		var err error
		pgClient, err = xray.NewClientWithConfig(config, types.Postgres)
		if err != nil {
			fmt.Printf("Error connecting to PostgreSQL: %v", err)
			os.Exit(1)
		}

		fmt.Println("Connected to PostgreSQL")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connected to PostgreSQL")
	},
}

var PgExecuteCmd = &cobra.Command{
	Use:    "execute",
	Short:  "Execute a SQL query in postgres",
	Args:   cobra.ExactArgs(1),
	PreRun: postgresCmd.PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		data, err := pgClient.Execute(query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return
		}
		fmt.Println(data)
	},
}

var PgSchemaCmd = &cobra.Command{
	Use:    "schema",
	Short:  "Extract the schema of a Postgres database",
	Args:   cobra.ExactArgs(1),
	PreRun: postgresCmd.PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		table := args[0]
		schema, err := pgClient.Schema(table)
		if err != nil {
			fmt.Println("Error getting schema:", err)
			return
		}
		fmt.Println(schema)
	},
}

var PgTablesCmd = &cobra.Command{
	Use:    "tables",
	Short:  "List all tables in the Postgres database",
	Args:   cobra.ExactArgs(1),
	PreRun: postgresCmd.PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		DBName := args[0]
		tables, err := pgClient.Tables(DBName)
		if err != nil {
			fmt.Println("Error getting tables:", err)
			return
		}
		fmt.Println(tables)
	},
}

func init() {
	flags := postgresCmd.PersistentFlags()
	flags.String("host", "", "PostgreSQL server host")
	flags.String("port", "5432", "PostgreSQL server port")
	flags.String("user", "", "PostgreSQL user")
	flags.String("database", "", "PostgreSQL database")
	flags.String("ssl", "", "SSL mode")

	postgresCmd.AddCommand(PgExecuteCmd)
	postgresCmd.AddCommand(PgSchemaCmd)
	postgresCmd.AddCommand(PgTablesCmd)
}
