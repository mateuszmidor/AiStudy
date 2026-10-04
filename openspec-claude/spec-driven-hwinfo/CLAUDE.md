# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`hwinfo` is a Go CLI that prints basic hardware information. Currently it only reports CPU info, and only on macOS (reads via sysctl, no cgo, no subprocesses).

## Commands

```
go build ./...
go vet ./...
go test ./...
go test ./internal/cpu -run TestVendorFromModel   # single test
go run .                                          # print the CPU section
```

## Architecture

`main.go` only wires `cpu.SystemSource{}` into `cpu.Read` and prints `cpu.Format(...)`. All logic lives in `internal/cpu`, split so it is testable on any machine:

- `cpu.go` — the `Source` interface (`String`/`Uint32`/`Uint64` lookups by sysctl name, returning an error when the key doesn't exist), the `CPU` struct, and `Read(Source)`.
- `format.go` — `Format(CPU)` renders the labeled text section.
- `sysctl_darwin.go` — `SystemSource`, the real `Source` backed by `golang.org/x/sys/unix`. It is the only file behind `//go:build darwin`, so on non-darwin platforms `cpu.SystemSource` doesn't exist and `main.go` won't build (other platforms are an explicit non-goal).
- Tests use a `fakeSource` (maps of canned values) instead of real sysctl, which is how Intel-style brand strings and missing-frequency cases are covered on an Apple Silicon machine.

Conventions that span files:
- A zero value in `CPU` means "unavailable" (no separate availability flags). `Read` leaves a field zero when its lookup fails; only `Format` turns that into the text `unavailable`. A missing sysctl key is a normal condition, never an error, and `main` always exits 0.
- Cores are reported as unavailable if either the physical or logical count is missing.
- On Apple Silicon `machdep.cpu.vendor` and `hw.cpufrequency_max` don't exist, so vendor is derived from the first word of the brand string (trademark marks `(R)`/`(TM)` stripped) and speed shows `unavailable`.

## Spec-driven workflow (OpenSpec)

Changes are planned with OpenSpec (`openspec/config.yaml`, schema `spec-driven`). Each change lives in `openspec/changes/<name>/` with `proposal.md`, `design.md`, `tasks.md` and `specs/`; the existing `add-cpu-info` change holds the design rationale for the code above. Use the `/opsx:*` commands (`propose`, `apply`, `verify`, `archive`, …) or the matching `openspec-*` skills in `.claude/` rather than editing artifacts ad hoc.
