## 1. Project Setup

- [x] 1.1 Create main.go with basic HTTP server structure
- [x] 1.2 Create ./data directory on startup
- [x] 1.3 Configure server to listen on localhost:8080

## 2. Backend - Data Model and Storage

- [x] 2.1 Define Note struct (id, title, content, createdAt, updatedAt, tags)
- [x] 2.2 Implement unique ID generation function
- [x] 2.3 Implement note save function (JSON file in ./data)
- [x] 2.4 Implement note load function (read all notes from ./data)
- [x] 2.5 Implement data directory creation if not exists

## 3. Backend - HTTP Handlers

- [x] 3.1 Create GET / handler for main board page
- [x] 3.2 Create GET /notes handler returning HTML with all notes
- [x] 3.3 Create POST /notes handler for note creation
- [x] 3.4 Create GET /notes/:id handler for single note view
- [x] 3.5 Add basic input validation (title and content required)

## 4. Frontend - HTML Structure

- [x] 4.1 Create index.html with board layout
- [x] 4.2 Add Bootstrap 5 CDN link
- [x] 4.3 Create Add button
- [x] 4.4 Create note creation modal form
- [x] 4.5 Create note view modal dialog

## 5. Frontend - JavaScript

- [x] 5.1 Implement note creation form submission
- [x] 5.2 Implement note card click to open view modal
- [x] 5.3 Implement board refresh after note creation
- [x] 5.4 Implement tag parsing (comma-separated, trim whitespace)
- [x] 5.5 Add error handling for form validation

## 6. Static Assets

- [x] 6.1 Create custom CSS for note card styling
- [x] 6.2 Configure static file serving for /static route

## 7. HTML Templates

- [x] 7.1 Create main page template with board grid
- [x] 7.2 Create note card template
- [x] 7.3 Create create modal template
- [x] 7.4 Create view modal template