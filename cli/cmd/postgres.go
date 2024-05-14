/*
Copyright © 2024 Xray
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	"github.com/thesaas-company/xray/types"
)

var pgClient types.ISQL

// postgresCmd represents the postgres command
var postgresCmd *cobra.Command

func init() {
	postgresCmd = &cobra.Command{
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
			fmt.Println("Hello from PostgreSQL shell!")
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("> ")
				cmdString, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				err = PostgresRun(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		},
	}
	flags := postgresCmd.PersistentFlags()
	flags.String("host", "", "PostgreSQL server host")
	flags.String("port", "5432", "PostgreSQL server port")
	flags.String("user", "", "PostgreSQL user")
	flags.String("database", "", "PostgreSQL database")
	flags.String("ssl", "", "SSL mode")

	// postgresCmd.AddCommand(PgExecuteCmd)
	// postgresCmd.AddCommand(PgSchemaCmd)
	// postgresCmd.AddCommand(PgTablesCmd)
}

func PostgresRun(cmdString string) (err error) {
	cmdString = strings.TrimSuffix(cmdString, "\n")
	arrCommandStr := strings.Fields(cmdString)

	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "execute":
		if len(arrCommandStr) < 2 {
			fmt.Println("Please provide a query to execute.")
			return
		}
		query := strings.Join(arrCommandStr[1:], " ")
		res, err := pgClient.Execute(query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return err
		}
		fmt.Println(res)
	case "tables":
		database, _ := postgresCmd.Flags().GetString("database")
		tables, err := pgClient.Tables(database)
		if err != nil {
			fmt.Println("Error fetching tables:", err)
			return err
		}
		fmt.Println(tables)
	case "schema":
		if len(arrCommandStr) < 2 {
			fmt.Println("Please provide a table name to get its schema.")
			return
		}
		table := arrCommandStr[1]
		schema, err := pgClient.Schema(table)
		if err != nil {
			fmt.Println("Error fetching schema:", err)
			return err
		}
		fmt.Println(schema)
	default:
		fmt.Println("Unknown command:", arrCommandStr[0])
	}

	return
}

// var PgExecuteCmd = &cobra.Command{
// 	Use:    "execute",
// 	Short:  "Execute a SQL query in postgres",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: postgresCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		query := args[0]
// 		data, err := pgClient.Execute(query)
// 		if err != nil {
// 			fmt.Println("Error executing query:", err)
// 			return
// 		}
// 		fmt.Println(data)
// 	},
// }

// var PgSchemaCmd = &cobra.Command{
// 	Use:    "schema",
// 	Short:  "Extract the schema of a Postgres database",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: postgresCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		table := args[0]
// 		schema, err := pgClient.Schema(table)
// 		if err != nil {
// 			fmt.Println("Error getting schema:", err)
// 			return
// 		}
// 		fmt.Println(schema)
// 	},
// }

// var PgTablesCmd = &cobra.Command{
// 	Use:    "tables",
// 	Short:  "List all tables in the Postgres database",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: postgresCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		DBName := args[0]
// 		tables, err := pgClient.Tables(DBName)
// 		if err != nil {
// 			fmt.Println("Error getting tables:", err)
// 			return
// 		}
// 		fmt.Println(tables)
// 	},
// }
