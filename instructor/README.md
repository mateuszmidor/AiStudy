# Random City with Instructor Go Demo

Asks OpenAI GPT-4o (via **Instructor Go**)
for one random city and prints it as `City { city, country }` typed JSON.

How is the valid JSON result achieved:  
you define a plain Go `struct`, and Instructor builds the JSON schema from it,
injects it into the prompt, parses and validates the model's answer, and retries
on failure. No server, no codegen, no manual JSON parsing.

## Run

```sh
export OPENAI_API_KEY=<your APIKEY>
go run .
```

```
2026/09/15 20:10:00 Random city: Chengdu, China
```

## Prerequisites

- Go 1.26+
- `OPENAI_API_KEY` (account with credit)
