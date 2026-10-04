## Why

Users need a simple, local-only web application for creating and viewing short notes. The current landscape lacks lightweight, self-contained note-taking tools that run locally without requiring databases or external dependencies. This WebNotesApp provides a zero-setup solution for personal note management.

## What Changes

- Create a Golang web application running on localhost:8080
- Implement a single-screen board-style UI with note cards in a responsive grid
- Add note creation via modal form with title, content, and optional tags
- Add note viewing via read-only modal dialog
- Store each note as a pretty-printed JSON file in ./data directory
- Implement static asset serving for HTML, JS, and Bootstrap 5

## Capabilities

### New Capabilities

- `note-creation`: Allow users to create notes with title, content, and optional tags via a modal form
- `note-display`: Display notes as cards in a responsive board layout ordered by newest first
- `note-viewing`: View full note details in a read-only modal dialog
- `data-persistence`: Store notes as individual JSON files with generated UIDs

### Modified Capabilities

None - this is a new capability with no existing specs.

## Impact

- New Go backend server at localhost:8080
- New HTML frontend with Bootstrap 5 UI
- New ./data directory for note storage
- Static assets (HTML, JS, CSS) served from disk