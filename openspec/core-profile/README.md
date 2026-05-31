# Core profile

OpenSpec initialized with the **core-profile** only installs the minimal OpenCode integration: four slash commands and their paired skills. Initialization:
```sh
openspec init
```

Command definitions live under [`.opencode/commands/`](./.opencode/commands/).

## Available commands

| Command | Description |
|---------|-------------|
| `/opsx:explore` | Enter explore mode — think through ideas, investigate problems, and clarify requirements (no implementation). |
| `/opsx:propose` | Propose a new change — create the change directory and generate proposal, design, and tasks artifacts in one step. |
| `/opsx:apply` | Implement tasks from an OpenSpec change. |
| `/opsx:archive` | Archive a completed change. |

Typical flow: **explore** (optional) → **propose** → **apply** → **archive**.

For the full OpenSpec tool install, workflows, and upgrading, see [../README.md](../README.md).
