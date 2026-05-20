## ADDED Requirements

### Requirement: Notes displayed in board grid layout
The system SHALL display notes as rectangular cards in a responsive grid layout that uses the available viewport.

#### Scenario: Grid layout on desktop
- **WHEN** user views the board on a desktop screen
- **THEN** notes are displayed in a multi-column grid

#### Scenario: Grid layout on mobile
- **WHEN** user views the board on a mobile screen
- **THEN** notes are displayed in a single column layout

### Requirement: Notes ordered by newest first
The system SHALL display notes sorted by createdAt timestamp in descending order.

#### Scenario: Notes sorted correctly
- **WHEN** multiple notes exist in the system
- **THEN** notes are displayed with newest at the top

### Requirement: Note card shows summary information
The system SHALL display each note card with title, created date, tags, and a short preview of content.

#### Scenario: Card displays all fields
- **WHEN** a note is displayed on the board
- **THEN** the card shows title, created date (formatted), tags (if any), and content preview (truncated)

### Requirement: Empty state shows only Add button
The system SHALL display only the Add button when no notes exist.

#### Scenario: Empty state
- **WHEN** no notes exist in the system
- **THEN** only the Add button is visible with no additional helper text