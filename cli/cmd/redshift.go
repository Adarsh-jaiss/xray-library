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

var RedshiftClient types.ISQL

// redshiftCmd represents the redshift command
var redshiftCmd = &cobra.Command{
	Use:   "redshift",
	Short: "Interact with Redshift databases",
	Long:  `This command allows you to interact with Redshift databases. You can use this command to connect to a Redshift database and run queries.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		accessKey, _ := cmd.Flags().GetString("access-key")
		secretAccessKey, _ := cmd.Flags().GetString("secret-access-key")
		DatabaseName, _ := cmd.Flags().GetString("database")
		Schema, _ := cmd.Flags().GetString("schema")

		config := config.Config{
			AWS: config.AWS{
				Region:          region,
				AccessKey:       accessKey,
				SecretAccessKey: secretAccessKey,
			},
			DatabaseName: DatabaseName,
			Schema:       Schema,
		}

		var err error
		RedshiftClient, err = xray.NewClientWithConfig(&config, types.Redshift)
		if err != nil {
			fmt.Printf("Error connecting to Redshift: %v", err)
			os.Exit(1)
		}

		fmt.Println("Connected to Redshift")
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
			err = RedshiftRunCommand(cmdString)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	},
}

func RedshiftRunCommand(cmdString string) (err error) {
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
		res, err := RedshiftClient.Execute(query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return err
		}
		fmt.Println(res)
	case "tables":
		database, _ := bigqueryCmd.Flags().GetString("database")
		tables, err := RedshiftClient.Tables(database)
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
		schema, err := RedshiftClient.Schema(table)
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

func init() {
	flag := redshiftCmd.PersistentFlags()
	flag.String("region", "", "AWS region")
	flag.String("access-key", "", "AWS access key")
	flag.String("secret-access-key", "", "AWS secret")
	flag.String("database", "", "Database name")
	flag.String("schema", "", "Schema name")

}
