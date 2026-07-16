See: [.opencode/skills](.opencode/skills), [.opencode/agents](.opencode/agents)

# OpenCode

## Install
```sh
# install opencode
pacman -Sy opencode

# enable text copy+paste in opencode
pacman -Sy xclip

# install local MCP for reading PDFs
pacman -Sy python-pip python-pipx
pipx ensurepath
pipx install 'pdf-mcp[multicolumn]'
bash # open new shell to have pdf-mcp available
```

## Config
In `$HOME/.config/opencode/opencode.jsonc`:
```json
{
  "$schema":"https://opencode.ai/config.json",
  "lsp":true,
  "plugin":[
    "superpowers@git+https://github.com/obra/superpowers.git"
  ],
  "mcp":{
    "pdf-mcp":{
      "type":"local",
      "command":[
        "pdf-mcp"
      ]
    },
    "czyjesteldorado": {
      "type": "remote",
      "url": "https://czyjesteldorado.pl/_mcp"
    }
  }
}
```
The above config will enable opencode with:
- Language Server Protocol (gopls)
- additional set of agent skills (superpowers)
- reading PDF files (json-mcp)
- finding new jobs :)

## Use it
- `opencode --continue` - continue last session
- `opencode debug config` - print effective config - all openconfig.jsonc files merged together
- `opencode models` - list available provider/model
- `opencode pr 123` - fetch PR #123 from github and run opencode; needs gh cli installed and authorized in github
- `opencode run "what programming language is used in this project"` - run command and exit
- in opencode:
    - `/init` - create AGENTS.md
    - `/undo`, `/redo` - revert/restore last change - uses git so project must be git-controlled
    - `!ls` - run shell command
    - `ctrl+x e` - open editor to write long prompt
    - `ctrl+p` - open menu, e.g. eg to select llm model
    - `@main.go` - attach file contents to context
- `opencode.jsonc`
    - global: `~/.config/opencode/opencode.jsonc`
    - project: `./opencode.jsonc`
    - what can be configured in opencode.jsonc:
        - allowed tools like "write", "bash", "glob"
        - available providers and models like ollama:qwen3.7
        - preconfigured agents and commands (literally inline definitions)
        - instructions (paths to files with instructions like AGENTS.md or CLAUDE.md)
        - MCP servers
        - plugins
        - Language Server Protocols; `add "lsp": true` and gopls is available to opencode

## OpenSpec - AI-coding framework
See: [../openspec](../openspec/)

## Creating skill

https://agentskills.io/skill-creation/best-practices

- **name** must be lower-case, can not start or end with '-', can not include '--' (double hyphens)
- **description** must explain when exactly to use this skill, so agent will pick it in all relevant use cases
    - if agent can achieve the goal without the skill, it may not use it, so tell him: "ALWAYS use this skill instead of coming up with ad-hoc solution"
- **body**
    - body is loaded only if agent decides that this skill should be used, so doesn't clutter context at startup
        - this is called: progressive disclosure
    - body should be narrowed down to specific task; big body equals more tokens fighting for LLM attention
    - should describe procedure: sequence of steps to achieve goal
        - every step should end in validation instruction, e.g. "validate by running: go vet ./..." so the agent can self-correct early
    - should explicitly say what tool to use for given task, and not a selection of tools, just one specific tool
        - can even say what NOT to do, e.g. "never use wget, always use curl"
    - should describe/explain things that will help LLM to achieve the goal
        - should not describe/explain things that LLM already knows (e.g. what http request is)
    - should contain GOTCHAs list in case there are gotchas in project (e.g. in table Users column uid is alias for user id)

Optimize skill by examining the agent's execution traces.
    - parts that the agent reinvents every time could be implemented as scripts under ./scripts

Skill can contain additional subfolders
- /assets - for e.g. templates to format output
- /scripts - for executable scripts e.g. bash, python, js
    - script should support '--help' flag - agent will use it when it fails to run the script
    - script should not expect user to type in password or anything; can't be interactive
- /references - for loading contidionally extra information into context (save tokens)