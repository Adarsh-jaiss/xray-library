/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// redshiftCmd represents the redshift command
var redshiftCmd = &cobra.Command{
	Use:   "redshift",
	Short: "Interact with Redshift databases",
	Long: `This command allows you to interact with Redshift databases. You can use this command to connect to a Redshift database and run queries.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("redshift called")
	},
}

func init() {

}
