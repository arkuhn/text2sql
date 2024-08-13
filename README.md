
Small DX project for generating SQL queries with semantic text / natural language. Allows generating/running/reprompting for sql queries
with GPT4 or Llama 3.1. Introspects a postgres DB and passes relevant table schema in prompt automatically.


Notes:
  - OpenAI integration assumes OPENAI_API_KEY is set in env.
  - Llama3.1 integration assumes you're running [Ollama](https://github.com/ollama/ollama?tab=readme-ov-file#linux)


Usage:
```
make build
./text2sql --help
./text2sql set-default-connection "postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName"
./text2sql query "query episodes from season 1 that have 'a' in the title or air date on an even month" --using Episode --model llama
```

