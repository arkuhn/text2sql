package llm

import (
        "context"
        "fmt"
        "regexp"
        "strings"
        "os"

        "github.com/sashabaranov/go-openai"
        "github.com/arkuhn/text2sql/internal/config"
)

func extractSQL(text string) string {
        sqlPattern := regexp.MustCompile("```\n(.*?)```")
        matches := sqlPattern.FindStringSubmatch(text)
        if len(matches) > 1 {
                return strings.TrimSpace(matches[1])
        }
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
        // Placeholder for local Llama model integration
        // You would implement the actual Llama integration here
        return fmt.Sprintf(`SELECT * FROM "%s" WHERE condition = 'placeholder';`, tables[0]), nil
}
