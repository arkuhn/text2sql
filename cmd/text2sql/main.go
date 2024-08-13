package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/arkuhn/text2sql/internal/config"
	"github.com/arkuhn/text2sql/internal/db"
	"github.com/arkuhn/text2sql/internal/llm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "text2sql",
	Short: "A CLI tool to generate SQL queries from natural language",
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Generate and execute SQL queries",
	Run:   runQuery,
}

var setDefaultConnectionCmd = &cobra.Command{
	Use:   "set-default-connection",
	Short: "Set the default database connection string",
	Run:   runSetDefaultConnection,
}

var setDefaultModelCmd = &cobra.Command{
	Use:   "set-default-model",
	Short: "Set the default LLM model",
	Run:   runSetDefaultModel,
}

func init() {
	rootCmd.AddCommand(queryCmd, setDefaultConnectionCmd, setDefaultModelCmd)

	queryCmd.Flags().StringSliceP("using", "u", nil, "Specify table names to use")
	queryCmd.Flags().StringP("connection", "c", "", "Database connection string")
	queryCmd.Flags().StringP("model", "m", "openai", "LLM model to use (llama, openai, claude)")
}

func runQuery(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Query text is required.")
		os.Exit(1)
	}
	queryText := args[0]

	using, _ := cmd.Flags().GetStringSlice("using")
	connection, _ := cmd.Flags().GetString("connection")
	model, _ := cmd.Flags().GetString("model")

	if connection == "" {
		connection = config.GetConfig("default_connection")
	}
	if model == "" {
		model = config.GetConfig("default_model")
	}

	if connection == "" {
		fmt.Println("Error: No connection string provided or found via config.")
		fmt.Println("Please set one via set-default-connection or provide one with --connection.")
		os.Exit(1)
	}

	var tables []string
	var err error
	if len(using) == 0 {
		tables, err = db.GetTableNames(connection)
	} else {
		tables = using
	}

	if err != nil || len(tables) == 0 {
		fmt.Println("Error: Unable to retrieve table names. Please check your connection string.")
		os.Exit(1)
	}

	schemas := make(map[string][]string)
	for _, table := range tables {
		schema, err := db.GetTableSchema(connection, table)
		if err != nil {
			fmt.Printf("Error fetching schema for table %s: %v\n", table, err)
			continue
		}
		schemas[table] = schema
	}

	originalQuery := queryText
	sqlQuery := ""

	for {
		if sqlQuery == "" {
			sqlQuery, err = llm.GenerateSQL(queryText, tables, schemas, model)
			if err != nil {
				fmt.Printf("Error generating SQL: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("Generated SQL Query:")
		fmt.Println(sqlQuery)

		fmt.Print("Choose an action: [r]un, [e]dit, [q]uit: ")
		var action string
		fmt.Scanln(&action)

		switch strings.ToLower(action) {
		case "r":
			result, err := db.ExecuteQuery(connection, sqlQuery)
			if err != nil {
				fmt.Printf("Query execution failed: %v\n", err)
				break
			}
			if len(result) > 0 {
				db.PrettyPrintQueryResults(result)
			} else {
				fmt.Println("Query executed successfully, but returned no results.")
			}
			return
		case "e":
			fmt.Print("Enter your refinement request: ")
			var refinement string
			fmt.Scanln(&refinement)
			newPrompt := fmt.Sprintf(`Original request: %s
Previously generated SQL: %s
Refinement request: %s

Please generate an updated SQL query based on the original request and the refinement.`, originalQuery, sqlQuery, refinement)
			sqlQuery, err = llm.GenerateSQL(newPrompt, tables, schemas, model)
			if err != nil {
				fmt.Printf("Error generating SQL: %v\n", err)
			}
		case "q":
			return
		default:
			fmt.Println("Invalid action. Please choose 'r', 'e', or 'q'.")
		}
	}
}

func runSetDefaultConnection(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Connection string is required.")
		os.Exit(1)
	}
	connection := args[0]
	config.SetConfig("default_connection", connection)
	fmt.Printf("Default connection set to: %s\n", connection)
	fmt.Println("Testing connection...")
	tables, err := db.GetTableNames(connection)
	if err != nil {
		fmt.Printf("Unable to retrieve tables: %v\n", err)
	} else {
		fmt.Printf("Connection successful. Available tables: %s\n", strings.Join(tables, ", "))
	}
}

func runSetDefaultModel(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Model name is required.")
		os.Exit(1)
	}
	model := args[0]
	if model != "llama" && model != "openai" && model != "claude" {
		fmt.Println("Invalid model. Choose from 'llama', 'openai', or 'claude'.")
		os.Exit(1)
	}
	config.SetConfig("default_model", model)
	fmt.Printf("Default model set to: %s\n", model)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
