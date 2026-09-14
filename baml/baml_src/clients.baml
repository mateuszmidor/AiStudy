// LLM client for this example, mirroring golang-openapi-starter/clients.baml
client<llm> GPT4oMini {
  provider openai
  options {
    model "gpt-4o-mini"
    api_key env.OPENAI_API_KEY
  }
}