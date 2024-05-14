/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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

var BigqueryClient types.ISQL

// bigqueryCmd represents the bigquery command
var bigqueryCmd *cobra.Command

func init() {
	bigqueryCmd = &cobra.Command{
		Use:   "bigquery",
		Short: "Interact with BigQuery databases",
		Long:  `This command allows you to interact with BigQuery databases. You can use this command to connect to a BigQuery database and run queries.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			ProjectID, _ := cmd.Flags().GetString("project-id")
			JSONKeyPath, _ := cmd.Flags().GetString("json-key-path")
			DatabaseName, _ := cmd.Flags().GetString("database")

			config := &config.Config{
				ProjectID:    ProjectID,
				JSONKeyPath:  JSONKeyPath,
				DatabaseName: DatabaseName,
			}

			var err error
			BigqueryClient, err = xray.NewClientWithConfig(config, types.BigQuery)
			if err != nil {
				fmt.Printf("Error connecting to BigQuery: %v", err)
				os.Exit(1)
			}

			fmt.Println("Connected to BigQuery")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from bigquery shell!")
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("> ")
				cmdString, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				err = BigQueryRunCommand(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		},
	}

	flags := bigqueryCmd.PersistentFlags()
	flags.String("project-id", "", "Project ID for BigQuery")
	flags.String("json-key-path", "", "Path to JSON key file for BigQuery")
	flags.String("database", "", "Database name for BigQuery")
}

func BigQueryRunCommand(commandStr string) (err error) {
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

