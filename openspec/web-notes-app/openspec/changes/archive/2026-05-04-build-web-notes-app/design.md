## Context

WebNotesApp is a local-only Golang web application for creating and viewing short notes. The target user is a single user running the app on localhost:8080. The application uses no external dependencies - only Go standard library for the backend and plain HTML/JavaScript with Bootstrap 5 for the frontend.

Key constraints:
- Backend: Go stdlib only (net/http, html/template, encoding/json, os/ioutil)
- Frontend: Plain HTML, vanilla JS, Bootstrap 5 (no frameworks)
- Storage: Local filesystem (one JSON file per note in ./data directory)
- Runtime: localhost:8080

## Goals / Non-Goals

**Goals:**
- Create a functional note-taking web application
- Provide a clean board-style UI for viewing notes
- Enable note creation via modal form
- Enable note viewing via read-only modal
- Persist notes as JSON files on local filesystem

**Non-Goals:**
- User authentication/authorization
- Note editing or deletion
- Search or filtering
- Multi-page navigation
- Database or external storage

## Decisions

1. **JSON file storage over database** - Simple approach with no external dependencies. Each note is a separate JSON file named with a generated UID.

2. **Server-side HTML rendering with vanilla JS** - Using html/template for initial page render, with JavaScript handling dynamic modal interactions. This avoids frontend frameworks while maintaining clean code.

3. **Bootstrap 5 for UI** - Provides responsive layout, cards, forms, buttons, and modals out of the box without custom CSS overhead.

4. **Static assets on disk** - CSS, JS, and HTML served as separate files rather than embedded in the binary, making development and debugging easier.

5. **In-memory note loading** - Notes are loaded from disk on each request to ensure freshness, rather than caching in memory.

## Risks / Trade-offs

- [Risk] No data validation on backend - The frontend sends data directly to the backend which stores it. Mitigation: Add basic validation in the API handler.
- [Risk] File system race conditions - Multiple simultaneous creates could have issues. Mitigation: Single-user app makes this low risk.
- [Risk] No backup or data migration - Notes are simple JSON files that can be manually backed up. Mitigation: Acceptable for local-only personal use.