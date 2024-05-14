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

var snowflakeClient types.ISQL

// snowflakeCmd represents the snowflake command
var SnowflakeCmd *cobra.Command

func init() {
	SnowflakeCmd = &cobra.Command{
		Use:   "snowflake",
		Short: "Interact with Snowflake databases",
		Long:  `This command allows you to interact with Snowflake databases. You can use this command to connect to a Snowflake database and run queries.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			account, _ := cmd.Flags().GetString("account")
			user, _ := cmd.Flags().GetString("user")
			database, _ := cmd.Flags().GetString("database")
			warehouse, _ := cmd.Flags().GetString("warehouse")
			port, _ := cmd.Flags().GetString("port")
			schema, _ := cmd.Flags().GetString("schema")

			config := &config.Config{
				Account:      account,
				Username:     user,
				DatabaseName: database,
				Warehouse:    warehouse,
				Port:         port,
				Schema:       schema,
			}

			var err error
			snowflakeClient, err = xray.NewClientWithConfig(config, types.Snowflake)
			if err != nil {
				fmt.Printf("Error connecting to Snowflake: %v", err)
				os.Exit(1)
			}

			fmt.Println("Connected to Snowflake")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to Snowflake shell!")
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("> ")
				cmdString, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				err = SnowflakeRunCommand(cmdString)
				if err != nil {
					fmt.Println("CODE FATA!!!!")
					fmt.Fprintln(os.Stderr, err)
				}
			}
		},
	}

	flags := SnowflakeCmd.PersistentFlags()
	flags.String("account", "", "Snowflake account name")
	flags.String("user", "", "Snowflake username")
	flags.String("database", "", "Snowflake database name")
	flags.String("warehouse", "", "Snowflake warehouse name")
	flags.String("port", "", "Snowflake port")
	flags.String("schema", "", "Snowflake schema name")

	// Add snowflake command to the shell command
	// snowflakeCmd.AddCommand(snowflakeExecuteCmd)
	// snowflakeCmd.AddCommand(snowflakeTablesCmd)
	// snowflakeCmd.AddCommand(snowflakeSchemaCmd)

}

func SnowflakeRunCommand(cmdString string) (err error) {
	cmdString = strings.TrimSuffix(cmdString, "\n")
	arrCommandStr := strings.Fields(cmdString)

	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "execute":
		if len(arrCommandStr) < 2 {
			fmt.Println("Please provide a query to execute.")
			return nil
		}
		query := strings.Join(arrCommandStr[1:], " ")
		res, err := snowflakeClient.Execute(query)
		if err != nil {
			fmt.Println("Error executing query")
		}
		fmt.Println(res)
	case "tables":
		database, _ := SnowflakeCmd.Flags().GetString("database")
		tables, err := snowflakeClient.Tables(database)
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
		schema, err := snowflakeClient.Schema(table)
		if err != nil {
			fmt.Println("Error fetching schema")
		}
		fmt.Println(schema)
	default:
		fmt.Println("Unknown command:", arrCommandStr[0])
	}
	return
}

// var snowflakeExecuteCmd = &cobra.Command{
// 	Use:    "execute",
// 	Short:  "Execute a query in snowflake",
// 	Args:   cobra.ExactArgs(1),
// 	PreRun: snowflakeCmd.PreRun,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		query := args[0]
// 		res, err := snowflakeClient.Execute(query)
// 		if err != nil {
// 			fmt.Println("Error executing query:", err)
// 		}
// 		fmt.Println(res)
// 	},
// }

// var snowflakeTablesCmd = &cobra.Command{
// 	Use:    "tables",
// 	Short:  "List tables in the snowflake database",
// 	PreRun: snowflakeCmd.PreRun,
// 	Args:   cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		DBName := args[0]
// 		tables, err := snowflakeClient.Tables(DBName)
// 		if err != nil {
// 			fmt.Println("Error fetching tables:", err)
// 		}
// 		fmt.Println(tables)
// 	},
// }

// var snowflakeSchemaCmd = &cobra.Command{
// 	Use:    "schema",
// 	Short:  "Get the schema of a snowflake table",
// 	PreRun: snowflakeCmd.PreRun,
// 	Args:   cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		table := args[0]
// 		schema, err := snowflakeClient.Schema(table)
// 		if err != nil {
// 			fmt.Println("Error fetching schema:", err)
// 		}
// 		fmt.Println(schema)
// 	},
// }
