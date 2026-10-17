# Council Notes: proposal.md

## Author Summary

Primary agent authored the draft, because the council round could not run (see below). The proposal introduces one change-scoped capability, `ram-info`, that reports installed **total physical memory** as a single labeled line (`Total Memory: <n> GB`, bytes ÷ 1024³ truncated) and prints `Total Memory: N/A` at exit 0 when the read fails. It records macOS-only as a deliberate scope decision, lists the tool-level output conventions as living inside this capability for now, and states the deliberate divergence from the Linux twin `../hw`.

## Reviewer Challenges

- **The council round did not run.** `adversarial-author` failed with `OpenCode's free tier can only be used from within OpenCode...`, and a one-line probe of `adversarial-reviewer` failed identically. Neither agent file pins a model, so both inherit the default. No reviewer challenges exist for this artifact; there was no independent review.
- **No cross-model review means single-perspective bias is unmitigated.** Per the `adversarial-authoring` skill, this is the documented fallback path and the user chose it knowingly (option: "Primary-agent authoring, deviation logged"). The `proposal` rule in `openspec/config.yaml` ("Must use adversarial-authoring skill") is therefore only partly satisfied: the workflow ran, the council step could not.
- **Self-check substituted for review (primary agent, not independent).** Checks applied before writing: every template section present and in order; one capability named exactly `ram-info`, matching the change directory; modified capabilities empty and justified (`openspec/specs/` is empty); no implementation walkthrough beyond what Impact legitimately states; non-goals explicit; no project context or rules copied into the artifact; glossary terms used exactly and bolded in prose only.

## Resolutions

- Accepted: capability is change-scoped `ram-info` (user decision during the interview, overriding the tool-scoped wording in `decisions.md` D1). The rationale is recorded in the proposal's Capabilities section.
- Accepted: tool-level conventions (single labeled line, 1024³ **GB**, `N/A` + exit 0) live inside `ram-info` for now, with a note that a tool-level capability may be split out later.
- Accepted: macOS-by-intent is stated explicitly in "What Changes", so the corrected A2 reading is what survives into the artifacts.
- Accepted: deliberate twin divergence recorded in Impact, including that the twin's `add-ram-reporting` is open and never applied.
- Accepted: `Why` framed as product motivation only. The workflow-exercise framing (A3) is deliberately kept out of the artifact and recorded here instead.
- Deferred: the independent council review of this artifact. Re-run `adversarial-author` + `adversarial-reviewer` once subagent execution works, and treat its findings as a revision round over this file.
- Deferred: `decisions.md` O3 (test scope), O4 (`ram` API names), O5 (injection seam) — all three are design-phase questions and are named as unresolved in the proposal's Impact.

## Remaining Risks

- The artifact was written without independent review, so a reviewer may later find scope or clarity defects that the self-check above did not catch.
- Because the capability is change-scoped, the tool-level output conventions will be restated if a later metric capability lands. Accepted now, but it is a known duplication cost.
- **GB** meaning 1024³ reads as incorrect to an SI-minded reader; the glossary pins the meaning, and the label was chosen so the number matches what macOS itself displays.
- The read seam is undecided (O5). If the design concludes that `gopsutil` must sit behind an injected reader, `internal/ram` may need reshaping at apply time, which the proposal currently does not commit to either way.
