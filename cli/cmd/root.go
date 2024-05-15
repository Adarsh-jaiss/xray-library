package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thesaas-company/xray"
	"github.com/thesaas-company/xray/config"
	xrayTypes "github.com/thesaas-company/xray/types"
	"gopkg.in/yaml.v3"
)

var (
	verbose bool
	cfgFile string
	dbType  string
)

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Time    float64         `json:"time"`
	Error   string          `json:"error"`
}

type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers []string) *Table {
	return &Table{headers: headers}
}

func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

func (t *Table) String() string {
	// Find the maximum width of each column
	columnWidths := make([]int, len(t.headers))
	for i, header := range t.headers {
		columnWidths[i] = len(header)
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	// Create a format string based on the column widths
	var formatBuilder strings.Builder
	for _, width := range columnWidths {
		formatBuilder.WriteString(fmt.Sprintf("%%-%ds ", width))
	}
	formatString := formatBuilder.String()

	// Print the headers
	var result strings.Builder
	result.WriteString(fmt.Sprintf(formatString, toInterfaceSlice(t.headers)...))
	result.WriteRune('\n')

	// Print a separator line
	for _, width := range columnWidths {
		result.WriteString(strings.Repeat("-", width) + " ")
	}
	result.WriteRune('\n')

	// Print the rows
	for _, row := range t.rows {
		result.WriteString(fmt.Sprintf(formatString, toInterfaceSlice(row)...))
		result.WriteRune('\n')
	}

	return result.String()
}

func toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// Command for interacting with databases
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with databases",
	Run: func(cmd *cobra.Command, args []string) {
		if !verbose {
			logrus.SetOutput(ioutil.Discard)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
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
			fmt.Printf("Error: Failed to connect to database: %s: %v\n", dbType, err)
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

			var result QueryResult
			err = json.Unmarshal(b, &result)
			if err != nil {
				fmt.Println("Error parsing query result:", err)
				continue
			}

			if len(result.Rows) == 0 {
				fmt.Println("No results found.")
				continue
			}

			table := NewTable(result.Columns)
			for _, row := range result.Rows {
				stringRow := make([]string, len(row))
				for i, v := range row {
					switch value := v.(type) {
					case string:
						decodedValue, err := base64.StdEncoding.DecodeString(value)
						if err != nil {
							fmt.Println("Error decoding base64 value:", err)
							stringRow[i] = value // Use original value if decoding fails
						} else {
							stringRow[i] = string(decodedValue)
						}
					default:
						stringRow[i] = fmt.Sprintf("%v", value)
					}
				}
				table.AddRow(stringRow)
			}

			// Print the table
			fmt.Println(table.String())
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
	default:
		return xrayTypes.MySQL
	}
}
