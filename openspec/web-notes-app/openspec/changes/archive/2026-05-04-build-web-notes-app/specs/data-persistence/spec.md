## ADDED Requirements

### Requirement: Notes stored as JSON files
The system SHALL store each note as an individual pretty-printed JSON file in the ./data directory.

#### Scenario: Note file created
- **WHEN** a new note is created
- **THEN** a JSON file is created in ./data directory with the note's UID as filename

#### Scenario: Data directory created automatically
- **WHEN** application starts and ./data does not exist
- **THEN** the ./data directory is created automatically

### Requirement: Note JSON structure
The system SHALL store notes with id, title, content, createdAt, updatedAt, and tags fields.

#### Scenario: Note JSON format
- **WHEN** a note is saved
- **THEN** the JSON file contains: id (UID), title (string), content (string), createdAt (ISO8601), updatedAt (ISO8601), tags ([]string)

### Requirement: Unique ID generation
The system SHALL generate a unique identifier for each note.

#### Scenario: Each note has unique ID
- **WHEN** multiple notes are created
- **THEN** each note receives a different UID