---
name: find-jobs
description: Search, list, or compare IT job offers on the Polish market. Use whenever the user mentions jobs — searching for job offers, "znajdź oferty pracy", "ile płacą Golang developerowi", salary research, or comparing listings from job boards — even if they don't explicitly ask for a "job search". Fetches from the CzyJestEldorado and JustJoin.it MCP servers.
---

# Find Jobs Skill

Fetch IT job listings from a job board MCP server, filter them, and present them as a Markdown listing. Match the user's language in prose, but keep offer titles in their original language.

## Step 1 — Interpret intent

Extract from the user's message: technology, min salary, work mode, city, contract type, and any role to exclude. If the user gave **neither technology nor work mode**, ask one consolidated question before fetching — an unfiltered search returns mostly irrelevant offers. Otherwise proceed and state which defaults you applied (Step 3 lists them).

## Step 2 — Pick the server

- User named a specific server → use ONLY that one.
- Not named → ask the user to pick: **CzyJestEldorado** or **JustJoin.it**.
- User asked for **both** → run both fetches in the same turn (they are independent) and dedupe the final list by title + company, since the boards overlap.

## Step 3 — Build the query

Investigate all request parameters the server supports before calling it. Pass every filter the user gave, plus these defaults:

**Defaults** — mid+ seniority (intern/junior roles only if the user asks), and exclude roles that are mainly Frontend, Fullstack, QA, SRE, DevOps, or pure management — *unless the user explicitly asked for one of those*.

| Intent | CzyJestEldorado (`search_jobs`) | JustJoin.it (`find_jobs`) |
|---|---|---|
| Technology | `phrase` (use quotes for exact match) | `category` where it fits, else `keywords` |
| Exclude a term | `"-term"` inside `phrase` | not supported — filter after the fetch |
| Seniority | `seniority: ["mid","senior","c_level"]` — allowed: intern, junior, mid, senior, c_level | `experienceLevels: "mid,senior"` — allowed: intern, junior, mid, senior, manager, c_level |
| Work mode | `workModes` | `remoteWorkOptions` |
| Contract | `contractTypes` | `employmentTypes` (b2b, permanent, uoz, internship) |
| Min salary (PLN) | `minSalary` | `minSalary` |
| City | — (not supported) | `city` |
| Sort | `sortOrder: "newest" \| "best_paid" \| "recommended"` | `sortBy: "publishedAt" \| "salary"` |

In OpenCode, MCP tools are called through the `execute` (Code Mode) tool, e.g. `tools.czyjesteldorado.search_jobs({...})` — not directly. If the result is long and truncated, read the entire saved result file, from beginning to end, into the current context instead of delegating to subagents.

## Step 4 — Filter by the requested technology

Skip this step entirely if the user named no technology. For every offer:

- **KEEP** — the requested technology (including aliases, e.g. Go/Golang) is a main requirement: it appears in the title or among the listed required skills/keywords.
- **DISCARD** — it is mentioned only incidentally, as a bonus, or "nice to have".

Also discard any offer that is mainly for an excluded role (Frontend, Fullstack, QA, SRE, DevOps, management), unless the user asked for such a role.

## Step 5 — Format the output

Use this exact structure, an ordered list starting at 1, and `N/A` (never `-`) wherever a value is missing in the offer:

- Heading: `# <Technology> Job Offers <today's date, DD.MM.YYYY>`
- One `## <Server Name>` section per server queried
- Fields per offer, in this order:

| Field | Notes |
|---|---|
| Title | Original language, no translation |
| Company | As listed |
| Salary | PLN (gross/month), range if available; `N/A` if undisclosed |
| Mode | `remote` / `hybrid` / `office` |
| Location | City name(s); `N/A` if fully remote |
| Technologies | Comma-separated; the user's requested technology first |
| Contract | e.g. `B2B`, `UoP`; `N/A` if absent |
| Posted | YYYY-MM-DD; `N/A` if absent |
| Link | Clickable |

Example (date is illustrative — always use today's date):

# Golang Job Offers 17.08.2026

## MCP: CzyJestEldorado

1. **Mid/Senior Go Engineer with Web API experience**
   - Company: CodiLime
   - Salary: 17000-24000 PLN
   - Mode: remote
   - Location: N/A
   - Technologies: Go, Web API, SQL, NoSQL, Redis, MongoDB, Kafka, Kubernetes, OIDC, JWT
   - Contract: B2B
   - Posted: 2026-09-18
   - Link: https://czyjesteldorado.pl/praca/327710-mid-senior-go-engineer-with-web-api-experience-codilime

2. **Golang Developer**
   - Company: ITFS
   - Salary: 21800-25200 PLN
   - Mode: hybrid
   - Location: Gdańsk
   - Technologies: Golang, Kubernetes
   - Contract: UoP
   - Posted: N/A
   - Link: https://czyjesteldorado.pl/praca/326587-golang-developer-itfs
