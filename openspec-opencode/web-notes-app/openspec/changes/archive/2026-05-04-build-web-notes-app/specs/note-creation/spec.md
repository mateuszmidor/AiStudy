## ADDED Requirements

### Requirement: User can create a note with title and content
The system SHALL allow users to create a note by providing a title and content, where both fields are required.

#### Scenario: Successful note creation
- **WHEN** user fills in title and content fields and submits the create form
- **THEN** a new note is saved to the filesystem with generated UID, current timestamps, and provided data

#### Scenario: Missing title
- **WHEN** user leaves title empty and submits the create form
- **THEN** the form displays an error and note is not created

#### Scenario: Missing content
- **WHEN** user leaves content empty and submits the create form
- **THEN** the form displays an error and note is not created

### Requirement: User can add optional tags to a note
The system SHALL allow users to add zero or more tags to a note as a comma-separated string.

#### Scenario: Creating note with tags
- **WHEN** user enters "work, important" in tags field and submits
- **THEN** the note is saved with tags array ["work", "important"]

#### Scenario: Creating note with no tags
- **WHEN** user leaves tags field empty and submits
- **THEN** the note is saved with empty tags array []

#### Scenario: Tags with extra whitespace
- **WHEN** user enters "  work  ,  important  " in tags field
- **THEN** the note is saved with trimmed tags ["work", "important"]

### Requirement: Note creation returns updated board
The system SHALL return the updated note board HTML after successful note creation.

#### Scenario: Board refresh after creation
- **WHEN** user successfully creates a note
- **THEN** the modal closes and the board refreshes displaying the new note in the correct position