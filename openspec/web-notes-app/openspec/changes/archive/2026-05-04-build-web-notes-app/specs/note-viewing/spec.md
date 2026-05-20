## ADDED Requirements

### Requirement: User can view note details
The system SHALL display a read-only modal dialog when user clicks on a note card.

#### Scenario: View note modal opens
- **WHEN** user clicks on a note card
- **THEN** a modal dialog opens showing full note details

### Requirement: Modal displays all note fields
The system SHALL display title, created date, updated date, tags, and full plain-text content in the view modal.

#### Scenario: View modal shows all fields
- **WHEN** user opens a note's view modal
- **THEN** title, createdAt, updatedAt, tags, and full content are displayed

#### Scenario: Content preserves line breaks
- **WHEN** note content contains multiple lines
- **THEN** the view modal displays content with preserved line breaks

### Requirement: View modal is read-only
The system SHALL NOT allow editing of note fields in the view modal.

#### Scenario: No edit controls in view modal
- **WHEN** user views a note in the modal
- **THEN** all fields are displayed as text only with no input controls