# OpenSpec

https://github.com/Fission-AI/OpenSpec/
- lightweight when compared to github speckit

## Install openspec tool 

### Arch linux

```sh
sudo pacman -S nvm
echo 'source /usr/share/nvm/init-nvm.sh' >> ~/.bashrc  # or ~/.zshrc
exec $SHELL
nvm install --lts
nvm use --lts
npm install -g @fission-ai/openspec@latest
```

## Simple Workflow / openspec core profile

```/opsx-propose ──► /opsx-apply ──► /opsx-archive```

```sh
cd yourproject/
openspec init # initialize openspec for use in opencode
opencode # development happens inside opencode
# /opsx-explore <what-you-want-to-build> # optional step: explore ideas, ask questions, find alternatives before asking to propose a change
/opsx-propose <what-you-want-to-build> # propose a change - creates proposal.md, design.md, spec.md and tasks.md
/opsx-apply # implement all tasks from just created tasks.md
/opsx-archive # move generated md files to archive folder; development of this change is finished
```

## Advanced Workflow / openspec custom profile

```/opsx-new ──► /opsx-continue(proposal.md) ──► /opsx-continue(design.md) ──► /opsx-continue(spec.md) ──► /opsx-continue(tasks.md) ──►  /opsx-apply ──► /opsx-verify ──► /opsx-sync ──► /opsx-archive```

Note: after `/opsx:new`, you can `/opsx:ff` to create proposal.md+desing.md+spec.md+tasks.md in one shot.

```sh
cd yourproject/
openspec init # initialize openspec for use in opencode
openspec config profile 
# select: Delivery and Workflows
# then select: Both (skills+commands)
# then select: Explore, New, Continue, Apply, Fast-Forward, Sync, Archive, Verify 
opencode # development happens inside opencode
# /opsx:explore <what-you-want-to-build> # optional step: explore ideas, ask questions, find alternatives before asking to propose a change
/opsx-new <what-you-want-to-build> # propose a change - create proposal.md, design.md, spec.md and tasks.md
/opsx-continue # create proposal.md, you can verify it before proceeding
/opsx-continue # create desing.md, you can verify it
/opsx-continue # create spec.md, you can verify it
/opsx-continue # create tasks.md, you can verify it
/opsx-apply # implement all tasks from just created tasks.md
/opsx-verify # check if implementation is in line with specification
/opsx-sync # update ./openspec/specs/... with currently implemented change
/opsx-archive # move generated md files to archive folder; development of this change is finished
```

## Upgrading openspec tool

```sh
# after upgrading the openspec tool, run the update cmd in your project root dir:
openspec update # re-scans installed tools and regenerates all skills/commands
```

## OpenSpec on steroids - intent driven

This is an improvement for OpenSpec - enhanced schema and additional skills:
- description: https://intent-driven.dev/blog/
- ready-to-use leaven repo: https://github.com/intent-driven-dev/intent-driven-template