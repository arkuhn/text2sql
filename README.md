
Small DevEx application for semantic SQL query generation. Anecdotally, it tends to work best when you can prompt like a developer,
eg `query all <rough resemblence to tables> where <rough resemblance to column names> are <some kind of value>, but only for <conditions>`

Notes:

    - OpenAI integration assumes OPENAI_API_KEY is set in env.
    - Llama3.1 integration assumes you're running [Ollama]("https://github.com/ollama/ollama?tab=readme-ov-file#linux")


Usage:

    1. `make build`
    2. `./text2sql --help`    
    3. `./text2sql query "query episodes from season 1 that have 'a' in the title" --using Episode --model llama`

