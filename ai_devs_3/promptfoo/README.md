# promptfoo 

For prompt evaluation and testing. 

## Installation

```sh
npm install -g promptfoo
# OR
brew install promptfoo
```

## Environment

Set your OPENAI_API_KEY environment variable.

## Usage

```sh
promptfoo init # create yaml file
promptfoo eval --grader openai:gpt-4o-mini # run evaluation and display in console. Use gpt-4o--mini for llm-rubric assertions (default: gpt-4o)
promptfoo view # display evaluation results in browser
```

```
┌────────────────────────────────┬──────────────────────────────────────────────────────────┬──────────────────────────────────────────────────────────┐
│ subject                        │ [openai:gpt-4o-mini] List the main colors of             │ [openai:gpt-3.5-turbo] List the main colors of           │
│                                │ {{subject}}, format the result as JSON list, don't add   │ {{subject}}, format the result as JSON list, don't add   │
│                                │ any extra markers, decorators nor comments. Example      │ any extra markers, decorators nor comments. Example      │
│                                │ list: [color1, color2, color3]                           │ list: [color1, color2, color3]                           │
├────────────────────────────────┼──────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────┤
│ banana                         │ [PASS] ["yellow", "green", "brown"]                      │ [PASS] ["yellow", "green"]                               │
├────────────────────────────────┼──────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────┤
│ avocado                        │ [PASS] ["green", "dark green", "yellow", "brown"]        │ [PASS] ["green", "brown", "black"]                       │
├────────────────────────────────┼──────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────┤
│ forest                         │ [PASS] ["green", "brown", "yellow", "orange", "red",     │ [PASS] ["green", "brown", "black", "yellow", "orange",   │
│                                │ "blue", "gray"]                                          │ "red", "blue", "gray"]                                   │
├────────────────────────────────┼──────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────┤
│ ocean                          │ [PASS] ["blue", "green", "turquoise", "teal", "navy",    │ [PASS] ["blue", "green", "turquoise"]                    │
│                                │ "cyan", "aqua", "deep blue", "seafoam", "indigo"]        │                                                          │
├────────────────────────────────┼──────────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────┤
│ rainbow                        │ [PASS] ["Red", "Orange", "Yellow", "Green", "Blue",      │ [PASS] [                                                 │
│                                │ "Indigo", "Violet"]                                      │ "Red",                                                   │
│                                │                                                          │ "Orange",                                                │
│                                │                                                          │ "Yellow",                                                │
│                                │                                                          │ "Green",                                                 │
│                                │                                                          │ "Blue",                                                  │
│                                │                                                          │ "Indigo",                                                │
│                                │                                                          │ "Violet"                                                 │
│                                │                                                          │ ]                                                        │
└────────────────────────────────┴──────────────────────────────────────────────────────────┴──────────────────────────────────────────────────────────┘

```