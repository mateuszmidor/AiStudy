---
name: find-jobs
description: This skill is dedicated for finding and listing job offers in IT. Use when the user asks to find jobs.
---

# Find Jobs Skill

This skill fetches current job listings from the MCP server "czyjesteldorado", applies precise filters, and presents results as a Markdown table.

## Step-by-Step Instructions

### Step 1 — Fetch offers via MCP
In order to fetch current job offers for Golang dvelopers, contact MCP server "czyjesteldorado" with required phrases ["Go", "Golang"] and excluded phrases ["Frontend", "Fullstack", "QA", "SRE", "DevOps", "Manager"]; check with MCP for details on how to filter results.

### Step 2 — go through the offers one by one and filter by seniority
**Keep** offers where the declared seniority is one of:
- Mid / Regular
- Senior
- Lead / Principal / Staff / Architect / Expert

**Discard** offers declared as:
- Junior / Intern / Trainee / Entry-level

If seniority is not stated in the offer, **keep it**

### Step 3 — go through the offers one by one and filter by location
Filtering rules:
- For offers with **work mode** = **fully remote**  -> always KEEP
- For offers with **work mode** = **office** or **hybrid** and location one of ["Gdańsk", "Sopot", "Gdynia", "Trójmiasto"] -> KEEP
- Otherwise -> DISCARD

### Step 4 — Collect data per offer

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

### Step 5 - Format Output
Expected output format is a Markdown ordered list, where Salary = '-' means salary not provided in the offer.
Example output list:


# Golang Job Offers 15.04.2026

1. **Mid/Senior Go Engineer with Web API experience**
   - Company: CodiLime
   - Salary: 17000-24000 PLN
   - Mode: remote
   - Location: -
   - Technologies: Go, Web API, SQL, NoSQL, Redis, MongoDB, Kafka, Kubernetes, OIDC, JWT
   - Link: https://czyjesteldorado.pl/praca/327710-mid-senior-go-engineer-with-web-api-experience-codilime

2. **Golang Developer**
   - Company: ITFS
   - Salary: 21800-25200 PLN
   - Mode: hybrid
   - Location: Gdańsk
   - Technologies: Golang, Kubernetes
   - Link: https://czyjesteldorado.pl/praca/326587-golang-developer-itfs

3. **Senior Golang Developer**
   - Company: Atos Poland Global Services Sp. z o.o.
   - Salary: -
   - Mode: remote
   - Location: -
   - Technologies: Go, Kubernetes, GitOps, GitHub Actions, OpenShift
   - Link: https://czyjesteldorado.pl/praca/337707-senior-golang-developer-atos-poland-global-services-sp-z-o-o
