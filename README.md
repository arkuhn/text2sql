
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
❯ ./text2sql query "query episodes from season 1 that have 'a' in the title" --using Episode --model llama

Generated SQL Query:
SELECT *
FROM "Episode"
WHERE season = 1 AND title LIKE '%a%';

Choose an action: [r]un, [e]dit, [q]uit: e
Edit prompt: I meant the letter b, and limit to first two results.

Generated SQL Query:
SELECT *
FROM "Episode"
WHERE season = 1 AND title LIKE '%b%' LIMIT 2;

Choose an action: [r]un, [e]dit, [q]uit: r

UPDATEDAT                       SEASON  AIRDATE                         LOCATION                        CREATEDAT                       SYNOPSIS                        UPDATE                          NUMBER  TITLE           ISCLOSED        YELPURL
        ID                              RESTAURANTNAMEOLD       ISCLOSEDUPDATED                 RESTAURANTNAMENEW
2024-07-26 05:46:54.372 +0000   1       2011-07-31 00:00:00 +0000       3420 W Grace St, Chicago, IL    2024-07-26 05:46:54.372 +0000   Jon Taffer heads to Chicago,    Since the episode aired, The    3       Shabby Abbey    <nil>           https://www.yelp.com/biz/the-abbey-pub-chicago  clz2a5j5000029of9bow0j7e3       The Abbey Pub           <nil>
+0000                                   +0000                           60618                           +0000                           Illinois, to rescu...           Abbey Pub faced a ...

2024-07-26 05:53:55.976 +0000   1       2011-08-14 00:00:00 +0000       10 S Front St, Philadelphia,    2024-07-26 05:46:55.828 +0000   Jon Taffer heads to             Since the episode's airing,     5       Swanky Troubles <nil>           https://www.yelp.com/biz/sheer-philadelphia     clz2a5k9f00049of9ge1g3y4w       Swanky Bubbles          2024-07-26 05:53:55.974 +0000   Sheer
+0000                                   +0000                           PA 19106                        +0000                           Philadelphia, PA, to rescue...  Sheer has seen mixe...
                                                                +0000

Total rows: 2
```

