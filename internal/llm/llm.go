package llm

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/arkuhn/text2sql/internal/config"
	"github.com/sashabaranov/go-openai"
)

func extractSQL(text string) string {
	sqlPattern := regexp.MustCompile("(?s)```(?:sql)?\n?(.*?)\n?```")
	matches := sqlPattern.FindStringSubmatch(text)
	if len(matches) > 1 {
		extractedSQL := strings.TrimSpace(matches[1])
		return extractedSQL
	}

	//log.Println("No SQL code block found. Returning original text.")
	return strings.TrimSpace(text)
}

func quoteTableNames(sql string, tables []string) string {
	//for _, table := range tables {
	//        pattern := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(table)))
	//        sql = pattern.ReplaceAllString(sql, fmt.Sprintf(`"%s"`, table))
	//}
	return sql
}

func GenerateSQL(query string, tables []string, schemas map[string][]string, model string) (string, error) {
	var rawSQL string
	var err error

	switch model {
	case "llama":
		rawSQL, err = generateSQLLlama(query, tables, schemas)
	case "openai":
		rawSQL, err = generateSQLOpenAI(query, tables, schemas)
	case "claude":
		rawSQL, err = generateSQLClaude(query, tables, schemas)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}

	if err != nil {
		return "", err
	}

	extractedSQL := extractSQL(rawSQL)
	return quoteTableNames(extractedSQL, tables), nil
}

func generateSQLOpenAI(query string, tables []string, schemas map[string][]string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = config.GetConfig("OPENAI_API_KEY")
	}

	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not seti in env or config")
	}

	client := openai.NewClient(apiKey)
	schemaInfo := ""
	for table, columns := range schemas {
		schemaInfo += fmt.Sprintf("%s columns: %s\n", table, strings.Join(columns, ", "))
	}

	prompt := fmt.Sprintf(`Generate an SQL query for the following request: %s
Available tables and their schemas:
%s
IMPORTANT: Ensure to quote all table names in double quotes to preserve case sensitivity, e.g., "TableName".`, query, schemaInfo)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func generateSQLClaude(query string, tables []string, schemas map[string][]string) (string, error) {
	//apiKey := config.GetConfig("ANTHROPIC_API_KEY")
	//if apiKey == "" {
	//      return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	//}

	//client := anthropic.NewClient(apiKey)
	//schemaInfo := ""
	//for table, columns := range schemas {
	//      schemaInfo += fmt.Sprintf("%s columns: %s\n", table, strings.Join(columns, ", "))
	//}

	//prompt := fmt.Sprintf(`Generate an SQL query for the following request: %s
	//Available tables and their schemas:
	//%s
	//IMPORTANT: Ensure to quote all table names in double quotes to preserve case sensitivity, e.g., "TableName".`, query, schemaInfo)

	//resp, err := client.CompletionCreate(context.Background(), &anthropic.CompletionRequest{
	//      Model:      anthropic.Claude3Sonnet,
	//      Prompt:     prompt,
	//      MaxTokens:  300,
	//      Temperature: 0,
	//})

	//if err != nil {
	//      return "", err
	//}

	return fmt.Sprintf(`placeholder`), nil
}

func generateSQLLlama(query string, tables []string, schemas map[string][]string) (string, error) {
	// Check if ollama is installed
	if _, err := exec.LookPath("ollama"); err != nil {
		return "", fmt.Errorf("ollama is not installed. Please install it from https://ollama.ai")
	}

	// Check if ollama server is running
	if err := checkOllamaServer(); err != nil {
		// Start ollama server
		if err := startOllamaServer(); err != nil {
			return "", fmt.Errorf("failed to start ollama server: %v", err)
		}
	}

	// Prepare the prompt
	schemaInfo := ""
	for table, columns := range schemas {
		schemaInfo += fmt.Sprintf("%s columns: %s\n", table, strings.Join(columns, ", "))
	}

	prompt := fmt.Sprintf(`Generate an SQL query for the following request: %s
Available tables and their schemas:
%s
IMPORTANT: Ensure to quote all table names in double quotes to preserve case sensitivity, e.g., "TableName".
    Additionally, your output will be sent right to a database so output no information other than the query, 
    otherwise it will fail when sent to the server.`, query, schemaInfo)

	// Run ollama
	cmd := exec.Command("ollama", "run", "llama3.1", prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run ollama: %v", err)
	}

	return out.String(), nil
}

func checkOllamaServer() error {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama server is not running")
	}
	return nil
}

func startOllamaServer() error {
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait for the server to start
	for i := 0; i < 10; i++ {
		if err := checkOllamaServer(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("ollama server failed to start within the expected time")
}
