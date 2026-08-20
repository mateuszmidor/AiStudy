---
name: find-jobs
description: Use when the user asks to find or list IT job offers, or asks to search for jobs.
---

# Find Jobs Skill

Fetch IT job listings from MCP servers and present them as a Markdown list.

**CORE PRINCIPLE: Process each MCP server separately and sequentially.** Run Steps 2-5 to completion for ONE server before touching the next server. NEVER fetch from multiple servers first and merge the raw results before filtering.

**No exceptions:**
- Don't fetch from several servers in parallel, even if the calls seem independent
- Don't merge filtered results of different servers into one renumbered list
- Don't skip the filtering for a server just because you will "filter later on the combined results"
- Each server's result must keep its own separate section in the output

**Red flags — STOP and redo:**
- Issuing a fetch call to a second MCP server before the first server's Steps 2-5 are all finished
- "The fetches are independent, so parallel is fine"
- "I'll merge the survivors after filtering each server"
- One combined list without per-server sections

## Step-by-Step Instructions

### Step 1 — Determine the MCP servers to be used

- If the user named a specific MCP server → use ONLY that server.
- Otherwise → use ALL available MCP servers that provide job offers. Enumerate the available MCP servers and their tools/resources; include only offer-providing servers (e.g. czyjesteldorado, justjoinit) and skip non-offer servers (e.g. pdf-mcp).

Then process each selected server through Steps 2-5 IN SEQUENCE, one server at a time. Do not fetch from all servers in parallel and merge.

### Step 2 — Fetch offers via MCP (current server)

Contact the server to fetch the offers. IMPORTANT: 
* make sure to only search for non-entry level roles - ask the MCP for mid/senior/principal/lead/architect roles
* pass any user-provided parameters that will narrow down the search (e.g. technology, min salary, work mode, city).
* double check with the MCP for details on supported filtering parameters before asking the MCP for job offers.
* if the result returned from MCP is long and gets truncated - read the entire result from file (from beginning to the end) into the current context without delegating to subagents.

### Step 3 — Filter by location (current server's results)

- If offer has **work mode** = **fully remote** -> always KEEP
- If offer has **work mode** = **office** or **hybrid** and location one of ["Gdańsk", "Sopot", "Gdynia", "Trójmiasto"] -> always KEEP
- In all other cases -> DISCARD

### Step 4 — Collect data per offer (current server's results)

For each remaining offer, collect:

| Field | Notes |
|---|---|
| Title | Original language, no translation |
| Company | As listed |
| Salary | In PLN (gross/month), range if available; `-` if not provided |
| Work Mode | `remote` / `hybrid` / `office` |
| Location | City name(s); `-` if fully remote |
| Technologies | Comma-separated, highlight Go/Golang first |
| Link | Clickable markdown link |

### Step 5 — Format output for this server

For THIS server, produce its own Markdown section with a heading `## <Server Name>` and a separate ordered list starting at 1 (numbering restarts at 1 inside each `##` section; never continue numbering across sections); use 'N/A' where salary or location is not provided in the offer. Once this server's list is complete, move to the next server and repeat Steps 2-5. The final answer contains one `##` section per server, in the order processed. Never merge the per-server sections into a single renumbered list (duplicates across servers are fine, keep them in their own sections).

Example output:

# Golang Job Offers 17.08.2026

## CzyJestEldorado

1. **Mid/Senior Go Engineer with Web API experience**
   - Company: CodiLime
   - Salary: 17000-24000 PLN
   - Mode: remote
   - Location: N/A
   - Technologies: Go, Web API, SQL, NoSQL, Redis, MongoDB, Kafka, Kubernetes, OIDC, JWT
   - Link: https://czyjesteldorado.pl/praca/327710-mid-senior-go-engineer-with-web-api-experience-codilime

2. **Golang Developer**
   - Company: ITFS
   - Salary: 21800-25200 PLN
   - Mode: hybrid
   - Location: Gdańsk
   - Technologies: Golang, Kubernetes
   - Link: https://czyjesteldorado.pl/praca/326587-golang-developer-itfs

## JustJoin.it

1. **Senior Golang Developer**
   - Company: Atos Poland Global Services Sp. z o.o.
   - Salary: N/A
   - Mode: remote
   - Location: N/A
   - Technologies: Go, Kubernetes, GitOps, GitHub Actions, OpenShift
   - Link: https://justjoin.it/job-offer/337707-senior-golang-developer