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

func Execute() {
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
	shellCmd.AddCommand(bigqueryCmd)
	shellCmd.AddCommand(redshiftCmd)
	// shellCmd.AddCommand(SnowflakeCmd)
	// Add subcommands to the serve command

}
