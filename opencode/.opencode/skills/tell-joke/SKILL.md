---
name: tell-joke
description: This skill returns a single random joke.
---

# Tell Joke Skill

This skill fetches single joke, then formats it according to template, then returns it to the caller.

## Step-by-Step Instructions

### Step 1 - Fetch the joke

Execute script ./scripts/fetch_joke.sh to get a single random joke.

### Step 2 - format the joke text according to template

The joke fetched in Step 1 must be formatted according to ./assets/output_template.md

### Step 3 - return the joke to caller

The formatted joke must be returned to the caller.