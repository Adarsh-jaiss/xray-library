package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
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
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			if input == "exit" {
				fmt.Println("Exiting shell.")
				break
			}
			fmt.Println(input)
		}
	},
}


func init() {
	// Add subcommands to the shell command
}

func main() {
	rootCmd := &cobra.Command{Use: "xray"}
	rootCmd.AddCommand(shellCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
