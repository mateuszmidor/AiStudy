# WebNotesApp

## Initialize openspec

```sh
openspec init
```

## Use LLM to create app description
```
I want to create a golang web-app WebNotesApp that will provide CRUD operations for short notes, with html interface, but running locally only (I'm the only user). Ask me all clarification questions needed to design the app, one by one, and finally draft a description of the app that will be the input to design phase. Go!
```

Generated App Description:

```
WebNotesApp is a local-only Golang web application for creating and viewing short notes through an HTML interface, running on localhost:8080. It is intended for a single user and stores data on the local filesystem, not in a database.

The backend must use only the Go standard library, with no third-party Go packages. Expected backend building blocks include net/http for routing and HTTP serving, html/template for server-rendered HTML, encoding/json for note serialization, and filesystem packages from the standard library for reading and writing note files.

The frontend must use plain HTML, vanilla JavaScript, and Bootstrap 5 as the UI library, with no frontend frameworks. Bootstrap should provide the main UI primitives such as responsive layout, cards, forms, buttons, and modal dialogs.

The app stores each note as its own pretty-printed JSON file inside a ./data directory resolved relative to the working directory from which the application is started, and the app creates that folder automatically if it does not already exist. Each note file is named with a generated UID, while the note payload contains title, content, createdAt, updatedAt, and tags.

The UI is a single-screen board-style web app that uses the available viewport to present notes as rectangular cards in a responsive board/grid layout. There is no multi-page navigation, no search, no filtering, no editing, and no deletion in this first phase.

The primary action is an Add button that opens a modal form for creating a note. In that form, title and content are required, tags are optional, and tags are entered as a comma-separated string that is parsed into a []string with surrounding whitespace trimmed.

After successful creation, the modal closes and the board refreshes so the new card appears in the correct position. Notes are ordered by newest created first.

Each note card shows the note title, created date, tags, and a short preview of the content. Clicking a note opens a read-only modal dialog showing title, created date, updated date, tags, and full plain-text content with preserved line breaks.

Static assets should be kept as separate files on disk rather than embedded into the binary. The initial empty state contains only the Add button, with no extra helper text.
Constraints

    Backend: Go stdlib only.

    Frontend: HTML + vanilla JS + Bootstrap 5.

    No frontend frameworks.

    No database; one JSON file per note.

    Local-only runtime on localhost:8080.
```

## Generate change specification - run from opencode

```sh
/opsx-propose <App Description>
```

## Implement

```sh
/opsx-apply
```