See: ./opencode/skills, ./opencode/agent

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