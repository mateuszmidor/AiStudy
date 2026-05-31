# Custom profile

OpenSpec initialized with the **custom-profile** provides the customized OpenCode command set: **all commands available under** [`.opencode/commands/`](./.opencode/commands/). Initialization:
```sh
openspec init
openspec config profile # here we select the commands 
```

## Available commands

| Command | Description |
|---------|-------------|
| `/opsx:explore` | Enter explore mode — think through ideas, investigate problems, and clarify requirements (no implementation). |
| `/opsx:new` | Start a new change using the experimental artifact workflow (OPSX). |
| `/opsx:continue` | Continue working on a change — create the next artifact (experimental). |
| `/opsx:ff` | Fast-forward artifact creation — generate everything needed to start implementation. |
| `/opsx:apply` | Implement tasks from an OpenSpec change (experimental). |
| `/opsx:sync` | Sync delta specs from a change to main specs. |
| `/opsx:verify` | Verify implementation matches change artifacts before archiving. |
| `/opsx:archive` | Archive a completed change in the experimental workflow. |
| `/opsx:onboard` | Guided onboarding — walk through a complete OpenSpec workflow cycle with narration. |

For full OpenSpec tool install, profiles, and upgrading, see [../README.md](../README.md).

