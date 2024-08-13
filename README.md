
Small DX project for generating SQL queries with semantic text / natural language.


Notes:
  - OpenAI integration assumes OPENAI_API_KEY is set in env.
  - Llama3.1 integration assumes you're running [Ollama](https://github.com/ollama/ollama?tab=readme-ov-file#linux)


Usage:
```
make build
./text2sql --help
./text2sql query "query episodes from season 1 that have 'a' in the title or air date on an even month" --using Episode --model llama
```

