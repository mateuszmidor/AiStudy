# connect opencode with ollama qwen3.5:9b

1. create file `~/.config/opencode/opencode.json`
    ```json
    {
      "$schema": "https://opencode.ai/config.json",
      "model": "ollama/qwen3.5:9b",
      "provider": {
        "ollama": {
          "npm": "@ai-sdk/openai-compatible",
          "name": "Ollama",
          "options": {
            "baseURL": "http://localhost:11434/v1"
          },
          "models": {
            "qwen3.5:9b": {
              "name": "Qwen3.5",
              "tools": true
            }
          }
        }
      }
    }
    ```

2. Make sure to run ollama with bigger than default 4k context window - 4k won't suffice to run tools:
      ```sh
      OLLAMA_CONTEXT_LENGTH=32000 ollama serve # may first need to: systemctl stop ollama
      ```

3. select "Qwen3.5 Ollama" model in opencode.