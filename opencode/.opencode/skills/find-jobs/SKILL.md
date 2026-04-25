---
name: find-jobs
description: This skill is dedicated for finding and listing job offers in IT. Use when the user asks to find jobs.
---

# Find Jobs Skill

Fetches current job listings from the MCP server "czyjesteldorado", applies precise filters, and presents results as a Markdown table.

## Step-by-Step Instructions

### Step 1 — Fetch offers via MCP
In order to fetch current job offers for Golang dvelopers, contact MCP server "czyjesteldorado" with required phrases ["Go", "Golang"] and excluded phrases ["Frontend", "Fullstack", "QA", "SRE", "DevOps", "Manager"].

### Step 2 — go through the offers one by one and filter by seniority
**Keep** offers where the declared seniority is one of:
- Mid / Regular
- Senior
- Lead / Principal / Staff / Architect / Expert

**Discard** offers declared as:
- Junior / Intern / Trainee / Entry-level

If seniority is not stated in the offer, **keep it**

### Step 3 — go through the offers one by one and filter by location
- For offers with **work mode** = **fully remote**  -> always KEEP
- For offers with **work mode** = **office** or **hybrid** and location one of ["Gdańsk", "Sopot", "Gdynia", "Trójmiasto"] -> KEEP
- Otherwise -> DISCARD

### Step 4 — Collect data per offer

For each remaining offer, collect:
| Field | Notes |
|---|---|
| Title | Original language, no translation |
| Company | As listed |
| Technologies | Comma-separated, highlight Go/Golang first |
| Salary | In PLN (gross/month), range if available; `-` if not provided |
| Location | City name(s); `-` if fully remote |
| Work Mode | `remote` / `hybrid` / `office` |
| Link | Clickable markdown link |

### Step 5 - Format Output
Expected output format is a Markdown table, where first column is item ordinal number, Salary = '-' means salary not provided in the offer.
Example output table:


# Golang Job Offers 15.04.2026

| No. | Title | Company | Technologies | Salary (PLN) | Location | Work Mode | Link |
|---|-------|---------|--------------|--------------|----------|-----------|------|
| 1. | Senior Engineer (Go) | Ericsson | C, Golang, Java, Python, LTE/4G/5G | 20000-25000 | Kraków, Łódź | hybrid | [Link](https://czyjesteldorado.pl/praca/326686-senior-engineer-c-or-go-ericsson) |
| 2. | Senior Backend Engineer (Ruby and/or Go), Tenant Scale | GitLab | Backend, Ruby, Go, GitLab, Security, AI, Architecture | - | - | remote | [Link](https://czyjesteldorado.pl/praca/300805-senior-backend-engineer-ruby-and-or-go-tenant-scale-cells-infrastructure-gitlab) |
| 3. | DevOps Engineer | ConnectPoint | Kubernetes, Terraform, Ansible, CI/CD, Python, Golang, PostgreSQL, SQL Server, Azure | 15000-20000 | Warszawa | hybrid, remote | [Link](https://czyjesteldorado.pl/praca/315328-devops-engineer-connectpoint) |
| 4. | Golang Architect | ITEAMLY | Golang, SQL, Python, Java, Pub/Sub, Kafka, NewRelic, DataDog | 26000-36000 | Kraków | remote | [Link](https://czyjesteldorado.pl/praca/323872-golang-architect-iteamly-spolka-z-ograniczona-odpowiedzialnoscia) |
| 5. | Senior Software Engineer II (Golang, Order) | SPOTON POLAND | Golang, PHP, React, JavaScript, AWS, Terraform | 23800-29800 | Kraków | hybrid | [Link](https://czyjesteldorado.pl/praca/316861-senior-software-engineer-ii-golang-order-spoton-poland-spolka-z-ograniczona-odpowiedzialnoscia) |
