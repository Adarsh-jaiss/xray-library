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

var client types.ISQL

// mysqlCmd represents the mysql command which is a subcommand of shell command
var mysqlCmd *cobra.Command

func init() {
	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "Interact with MySQL databases",
		PreRun: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetString("port")
			user, _ := cmd.Flags().GetString("user")
			// password, _ := cmd.Flags().GetString("password")
			database, _ := cmd.Flags().GetString("database")
			ssl, _ := cmd.Flags().GetString("ssl")

			config := &config.Config{
				Host:         host,
				DatabaseName: database,
				Username:     user,
				Port:         port,
				SSL:          ssl,
			}

			var err error
			client, err = xray.NewClientWithConfig(config, types.MySQL)
			if err != nil {
				fmt.Printf("Error connecting to MySQL: %v", err)
				os.Exit(1)
			}

			fmt.Println("Connected to MySQL")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from MySQL shell!")
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("> ")
				cmdString, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				err = MysqlRun(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		},
	}

	flags := mysqlCmd.PersistentFlags()
	flags.String("host", "", "MySQL host")
	flags.String("port", "", "MySQL port")
	flags.String("user", "", "MySQL user")
	flags.String("database", "", "MySQL database")
	flags.String("ssl", "", "Use SSL")

	// if err := mysqlCmd.ParseFlags(os.Args); err != nil {
	// 	fmt.Printf("Error parsing flags: %v\n", err)
	// 	os.Exit(1)
	// }

	// flags to the mysql command
	// mysqlCmd.AddCommand(MySQLExecuteCmd)
	// mysqlCmd.AddCommand(MySQLSchemaCmd)
	// mysqlCmd.AddCommand(MySQLTablesCmd)

}

func MysqlRun(cmdString string) (err error) {
	cmdString = strings.TrimSuffix(cmdString, "\n")
	arrCommandStr := strings.Fields(cmdString)

	if len(arrCommandStr) == 0 {
		fmt.Println("No command provided.")
		return
	}

	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "execute":
		if len(arrCommandStr) < 2 {
			fmt.Println("Please provide a query to execute.")
			return
		}
		query := strings.Join(arrCommandStr[1:], " ")
		res, err := client.Execute(query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return err
		}
		fmt.Println(res)
	case "tables":
		database, _ := postgresCmd.Flags().GetString("database")
		tables, err := client.Tables(database)
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
		schema, err := client.Schema(table)
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

// var MySQLExecuteCmd = &cobra.Command{
// 	Use:    "execute",
// 	Short:  "Execute a SQL query",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: mysqlCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// TODO: Implement the functionality here
// 		query := args[0]
// 		data, err := client.Execute(query)
// 		if err != nil {
// 			fmt.Printf("Error executing query: %v", err)
// 		}
// 		fmt.Println(data)

// 	},
// }

// // SchemaCmd represents the Schema command
// var MySQLSchemaCmd = &cobra.Command{
// 	Use:    "schema",
// 	Short:  "Get the schema of a table",
// 	PreRun: mysqlCmd.PreRun,
// 	Args:   cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		Table := args[0]
// 		schema, err := client.Schema(Table)
// 		if err != nil {
// 			fmt.Printf("Error getting schema: %v", err)
// 		}
// 		fmt.Println(schema)
// 	},
// }

// // TablesCmd represents the Tables command
// var MySQLTablesCmd = &cobra.Command{
// 	Use:    "tables",
// 	Short:  "Get the list of tables in the database",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: mysqlCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		DBName := args[0]
// 		data, err := client.Tables(DBName)
// 		if err != nil {
// 			fmt.Printf("Error getting tables: %v", err)
// 		}
// 		fmt.Println(data)
// 	},
// }
