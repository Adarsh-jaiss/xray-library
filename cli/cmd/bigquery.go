/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

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
				err = RunCommand(cmdString)
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

