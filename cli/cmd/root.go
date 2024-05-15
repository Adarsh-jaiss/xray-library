package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"io/ioutil"
	"gopkg.in/yaml.v3"
	"github.com/thesaas-company/xray/config"
	xrayTypes "github.com/thesaas-company/xray/types"
	"github.com/thesaas-company/xray"
)

var (
	verbose bool
	cfgFile string
	dbType string
)

// Command for interacting with databases
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with databases",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Verbose mode activated.")
		}
		
		if cfgFile == "" {
			fmt.Println("Error: Configuration file path is missing. Please use the --config flag to specify the path to your configuration file.")
			return
		}

		// Read the YAML file
		configData, err := ioutil.ReadFile(cfgFile)
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
			fmt.Printf("Error Failed to connect database: %s: %v\n", dbType, err)
			return
		}

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
			b, err := db.Execute(query)
			if err != nil {
				fmt.Println("Error executing query:", err)
				continue
			}
			fmt.Println(string(b))
		}
	},
}

func Execute() {
	rootCmd := &cobra.Command{Use: "xray"}

	rootCmd.AddCommand(shellCmd)
	shellCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	shellCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config.yaml")
	shellCmd.PersistentFlags().StringVarP(&dbType, "type", "t", "mysql", "Database type like mysql, postgres, bigquery")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {

}


// ParseDbType parses a string and returns the corresponding DbType.
func parseDbType(s string) (xrayTypes.DbType) {
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
	default:
		return xrayTypes.MySQL
	}
}