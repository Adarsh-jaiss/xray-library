package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

)

// Command for interacting with databases
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to database shell!")

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			query, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading query:", err)
				continue
			}
			query = strings.TrimSpace(query)
			if query == "exit" {
				fmt.Println("Exiting shell.")
				break
			}
			// Use xray lib to run the query and print the output like mysql and postgres cli
			fmt.Println(query)
			
		}
	},
}

// Command for interacting with databases
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to database serve!")
		// TDOD: Add Connect RPC Server
	},
}



func Execute()  {
	rootCmd := &cobra.Command{Use: "xray"}
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Add subcommands to the shell command
	shellCmd.AddCommand(mysqlCmd)
	shellCmd.AddCommand(postgresCmd)
	shellCmd.AddCommand(snowflakeCmd)
	shellCmd.AddCommand(bigqueryCmd)
	shellCmd.AddCommand(redshiftCmd)
	// Add subcommands to the serve command
	
}

func RunCommand(commandStr string) (err error) {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)

	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "execute":
		if len(arrCommandStr) < 2 {
			fmt.Println("Please provide a query to execute.")
			return
		}
		query := strings.Join(arrCommandStr[1:], " ")
		res, err := BigqueryClient.Execute(query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return err
		}
		fmt.Println(res)
	case "tables":
		database, _ := bigqueryCmd.Flags().GetString("database")
		tables, err := BigqueryClient.Tables(database)
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
		schema, err := BigqueryClient.Schema(table)
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

