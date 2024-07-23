# Ollama - run LLMs locally

## Install

https://ollama.com/

## Run

```sh
ollama run llama3
```
### API

https://github.com/ollama/ollama/blob/main/docs/api.md

### Generate

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3",
  "stream": false,
  "prompt":"Second planet of solar system?"
}'
```

Response:
```json
{"model":"llama3","created_at":"2024-07-23T06:52:36.614178Z","response":"The second planet in our solar system is Venus!","done":true,"done_reason":"stop","context":[128006,882,128007,271,16041,11841,315,13238,1887,30,128009,128006,78191,128007,271,791,2132,11841,304,1057,13238,1887,374,50076,0],"total_duration":2171653347,"load_duration":33036759,"prompt_eval_count":16,"prompt_eval_duration":220889000,"eval_count":11,"eval_duration":1916184000}
```

### Chat

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3",
  "stream": true,
  "messages": [
    { "role": "user", "content": "Second planet of solar system?" }
  ]
}'
```

Response:
```json
{"model":"llama3","created_at":"2024-07-23T06:54:10.764115Z","message":{"role":"assistant","content":"The"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:10.989147Z","message":{"role":"assistant","content":" second"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:11.215196Z","message":{"role":"assistant","content":" planet"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:11.441201Z","message":{"role":"assistant","content":" in"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:11.664788Z","message":{"role":"assistant","content":" our"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:11.891048Z","message":{"role":"assistant","content":" solar"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:12.071561Z","message":{"role":"assistant","content":" system"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:12.255476Z","message":{"role":"assistant","content":" is"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:12.441925Z","message":{"role":"assistant","content":" Venus"},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:12.622019Z","message":{"role":"assistant","content":"."},"done":false}
{"model":"llama3","created_at":"2024-07-23T06:54:12.807585Z","message":{"role":"assistant","content":""},"done_reason":"stop","done":true,"total_duration":2303979570,"load_duration":33080791,"prompt_eval_count":16,"prompt_eval_duration":225924000,"eval_count":11,"eval_duration":2043393000}
```