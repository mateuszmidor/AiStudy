# Random City with BAML Demo (Go)

Demo based on [BoundaryML/golang-openapi-starter](https://github.com/BoundaryML/baml-examples/tree/main/golang-openapi-starter):
asks OpenAI GPT-4o-mini (via BAML) for one random city and prints it as
`City { city, country }` typed JSON.

## TL;DR — Development Workflow

Two terminals:

Terminal 1 — start the BAML dev server (regenerates Go code on `baml_src/*.baml` file changes):

```bash
make install-tools # first time only
make serve
```

Terminal 2 — build and run:

```bash
make run # requires "make serve" as it will talk to the server who acts as middle man: app - baml_server - OpenAI
```

Edit `.baml` files in `baml_src/` while the server runs, then
`make run` again. See `make help` for all targets.

## Production

Your Go app is always an HTTP client to the BAML server, so production still
runs a server — but the *compiled* one (`baml serve`), not the dev server:

```bash
make serve-prod      # baml serve --port 2024 --dotenv (no codegen / file watching)
make run             # or build + deploy the binary separately
```

Deploy steps:

1. Generate the client at **build time** (CI-friendly, no server needed):
   `make generate` (runs `npx @boundaryml/baml generate`).
2. Run `make serve-prod` as a separate service side-by-side with the app (e.g.
   two containers on the same network).
3. Point the app at the server's real address via the server URL in
   `baml.NewConfiguration()` inside `main.go`.

`make build` self-generates `baml_client/` server-free if missing, so it works
in CI without the dev server.

## Prerequisites

- Node.js/npx
- `openapi-generator-cli`:
  `npm install -g @openapitools/openapi-generator-cli`
- Go 1.27
- `OPENAI_API_KEY` (account with credit)

---

# How to Work with BAML — Step-by-Step Guide

## The core idea

BAML is not a library, it's a **description + a server**. You never call
OpenAI from Go. You declare an LLM function in `.baml` files, a local BAML dev
server turns that declaration into a real HTTP endpoint (and actually runs it),
and your Go code is just an HTTP client to that server. Change the `.baml`
file → server hot-reloads and regenerates the Go client → the Go code doesn't
change.

<details>

### Layout (the only files you author by hand)

```
baml_src/            # your BAML source files
  clients.baml       # LLM provider config
  city.baml          # classes + functions + tests
  generators.baml    # says "generate a Go client from these functions"
baml_client/         # GENERATED Go package (never hand-edit; gitignore it)
main.go              # your Go entrypoint, uses baml_client/
```

---

## Step 1 — Declare clients (`baml_src/clients.baml`)

```baml
client<llm> GPT4oMini {
  provider openai
  options {
    model "gpt-4o-mini"
    api_key env.OPENAI_API_KEY   // reads from your environment
  }
}
```

This is a named, reusable "connection" to a model.

## Step 2 — Declare functions + types (`baml_src/city.baml`)

```baml
class City {
  city string
  country string
}

function GetRandomCity() -> City {
  client GPT4oMini              // which connection to use
  prompt #"
    {{ _.role("user") }}
    Pick one random city...     // your actual User prompt
    {{ ctx.output_format }}     // macro: injects the JSON schema for City
  "#
}

test RandomCityTest {
  functions [GetRandomCity]
  args {}
}
```

`ctx.output_format` is the magic macro — it tells the LLM the exact JSON shape
to return. `test` blocks give you a playground harness (run them from the BAML
extension or CLI).

## Step 3 — Tell it to make a Go client (`baml_src/generators.baml`)

```baml
generator target {
  output_type "rest/openapi"
  output_dir "../"
  version "0.226.2"              // must match your @boundaryml/baml version
  on_generate "openapi-generator-cli generate -i openapi.yaml -g go -o . ..."
}
```

`rest/openapi` = "expose these functions as HTTP + OpenAPI", and `on_generate`
= "after building the schema, run openapi-generator to emit a Go package".

## Step 4 — Start the server (Terminal 1)

```bash
cd baml && OPENAI_API_KEY=... npx @boundaryml/baml@0.226.2 dev
```

Every time you save/edit a `.baml` file, this:

1. compiles/validates the DSL (errors appear here),
2. serves each function as `POST /call/{FunctionName}` on `localhost:2024`,
3. regenerates `openapi.yaml` + `baml_client/` via the `on_generate` hook.

Sanity check: `curl http://localhost:2024/_debug/ping` → `pong`, and you
should see `baml_client/` appear.


## Step 5 — Write Go (Terminal 2, `main.go`)

```go
import baml "baml/baml_client"

cfg := baml.NewConfiguration()
b := baml.NewAPIClient(cfg).DefaultAPI

req := baml.NewGetRandomCityRequest()
resp, _, err := b.GetRandomCity(context.Background()).GetRandomCityRequest(*req).Execute()
fmt.Println(resp.City, resp.Country)   // e.g. "Kyoto Japan"
```

The generated code's exact method name/conventions come from *your* `function`
name in Step 2 (`GetRandomCity` → `b.GetRandomCity(...)`). If you add a new
function, inspect `baml_client/api_default.go` to see its generated signature.

## Step 6 — Run

```bash
go mod tidy     # pulls validator.v2 etc. the generated client needs
go run main.go
# 2026/09/14 20:08:07 Random city: Kyoto, Japan
```

---

## Your daily edit loop

1. Edit a `.baml` file (add a function, change a prompt, tweak a class).
2. Watch Terminal 1 regenerate + recompile (fix any DSL errors there).
3. Check `baml_client/api_default.go` for the new method name if you added a
   function.
4. `go build ./...` then `go run main.go`.

You stay in Go-land for logic/plumbing; all LLM behavior lives in `baml_src/`.

One trap from this project: the dev server re-creates a boilerplate
`baml_client/test/` with placeholder imports (`GIT_USER_ID/...`) on every
regenerate — harmless for `go build`/`go run`, but don't run `go test ./...`
inside `baml_client/`.

## File reference

- `baml_src/city.baml` — `City` class + `GetRandomCity` function
- `baml_src/clients.baml` — `GPT4oMini` OpenAI client
- `baml_src/generators.baml` — OpenAPI → Go codegen config
- `main.go` — calls `GetRandomCity` via the generated client
- `baml_client/` — generated Go client (gitignored)

</details>