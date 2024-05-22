package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/peterh/liner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	xrayTypes "github.com/thesaas-company/xray/types"
	"gopkg.in/yaml.v3"
)

// Command line flags
var (
	verbose bool
	cfgFile string
	dbType  string
	query   string
)

// QueryResult represents the result of a database query.
type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Time    float64         `json:"time"`
	Error   string          `json:"error"`
}

// Command for interacting with databases
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with databases",
	Long: `
	This command provides an interactive shell to execute SQL queries on various types of databases. 
	It supports MySQL, PostgreSQL, MSSQL, Redshift, Bigquery and Snowflake. 
	To use this command, you need to provide a configuration file with the --config flag or -c flag,
	and a database type with the --type flag or -t flag. 

	The configuration file should be in YAML format and contain the necessary database connection parameters 
	such as host, username, database name, port, and SSL settings.

	You can also control the verbosity of the command's output with the --verbose or -v flag. 
	When the verbose mode is on, the command will log additional information about its operation.

	In the interactive shell, you can type SQL queries and press Enter to execute them. 
	The results will be displayed in the console. Type 'exit' to leave the shell`,

	Run: func(cmd *cobra.Command, args []string) {

		// Set up logging
		if !verbose {
			logrus.SetOutput(io.Discard)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if cfgFile == "" {
			fmt.Println("Error: Configuration file path is missing. Please use the --config flag to specify the path to your configuration file.")
			return
		}

		// Read the YAML file
		configData, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("Error: Failed to read YAML file: %v\n", err)
			return
		}
		var cfg config.Config
		err = yaml.Unmarshal(configData, &cfg)
		if err != nil {
			fmt.Printf("Error: Failed to unmarshal YAML: %v\n", err)
			return
		}

		db, err := xray.NewClientWithConfig(&cfg, parseDbType(dbType))
		if err != nil {
			fmt.Printf("Error: Failed to connect to database: %s: %v\n", dbType, err)
			return
		}

		fmt.Println("Welcome to database shell!")

		if len(query) > 0 {
			if err := queryExecute(query, db); err != nil {
				fmt.Println(err)
				return
			}
		} else {
			line := liner.NewLiner()
			defer line.Close()

			for {
				line.SetCtrlCAborts(true)

				query, err := line.Prompt("> ")
				if err != nil {
					if err == liner.ErrPromptAborted {
						fmt.Println("Exiting shell.")
					} else {
						fmt.Println("Error reading query:", err)
					}
					break
				}

				if query == "exit" {
					fmt.Println("Exiting shell.")
					break
				}

				if err := queryExecute(query, db); err != nil {
					fmt.Println("Error executing query:", err)
				}

				line.AppendHistory(query)
			}
		}

	},
}

func queryExecute(query string, db xrayTypes.ISQL) error {

	b, err := db.Execute(strings.TrimSpace(query))
	if err != nil {
		fmt.Println("Error executing query:", err)
		return fmt.Errorf("error executing query result: %s", err)
	}

	var result QueryResult
	err = json.Unmarshal(b, &result)
	if err != nil {
		return fmt.Errorf("error parsing query result: %s", err)
	}

	if len(result.Rows) == 0 {

		return fmt.Errorf("no results found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(result.Columns)
	for _, row := range result.Rows {
		stringRow := make([]string, len(row))
		for i, v := range row {
			stringRow[i] = fmt.Sprintf("%v", v)
		}

		table.Append(stringRow)
	}

	// Print the table
	table.Render()
	return nil
}

// Execute runs the command line interface.
func Execute() {
	rootCmd := &cobra.Command{Use: "xray"}

	rootCmd.AddCommand(shellCmd)
	shellCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	shellCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config.yaml")
	shellCmd.PersistentFlags().StringVarP(&dbType, "type", "t", "mysql", "Database type like mysql, postgres, bigquery")
	shellCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Database query")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
}

// ParseDbType parses a string and returns the corresponding DbType.
func parseDbType(s string) xrayTypes.DbType {
	switch strings.ToLower(s) {
	case "mysql":
		return xrayTypes.MySQL
	case "postgres":
		return xrayTypes.Postgres
	case "snowflake":
		return xrayTypes.Snowflake
	case "bigquery":
		return xrayTypes.BigQuery
	case "redshift":
		return xrayTypes.Redshift
	case "mssql":
		return xrayTypes.MSSQL
	default:
		return xrayTypes.MySQL
	}
}