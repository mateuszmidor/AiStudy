package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Tags      []string  `json:"tags"`
}

type NoteForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func ensureDataDir() error {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		return os.Mkdir("data", 0755)
	}
	return nil
}

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

func parseTags(tagStr string) []string {
	if tagStr == "" {
		return []string{}
	}
	var tags []string
	for _, t := range splitAndTrim(tagStr, ",") {
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func splitAndTrim(s, sep string) []string {
	parts := split(s, sep)
	var result []string
	for _, p := range parts {
		result = append(result, trim(p))
	}
	return result
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func saveNote(note *Note) error {
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}
	filename := filepath.Join("data", note.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

func loadAllNotes() ([]Note, error) {
	entries, err := os.ReadDir("data")
	if err != nil {
		return nil, err
	}

	var notes []Note
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join("data", entry.Name()))
		if err != nil {
			continue
		}
		var note Note
		if err := json.Unmarshal(data, &note); err != nil {
			continue
		}
		notes = append(notes, note)
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].CreatedAt.After(notes[j].CreatedAt)
	})

	return notes, nil
}

func loadNote(id string) (*Note, error) {
	filename := filepath.Join("data", id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := loadAllNotes()
	if err != nil {
		notes = []Note{}
	}
	templates.ExecuteTemplate(w, "index.html", notes)
}

func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := loadAllNotes()
	if err != nil {
		notes = []Note{}
	}
	templates.ExecuteTemplate(w, "notes.html", notes)
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseMultipartForm(10 << 20)
	title := r.FormValue("title")
	content := r.FormValue("content")
	tagsStr := r.FormValue("tags")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	note := Note{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      parseTags(tagsStr),
	}

	if err := saveNote(&note); err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		return
	}

	notes, _ := loadAllNotes()
	templates.ExecuteTemplate(w, "notes.html", notes)
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := loadNote(id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	templates.ExecuteTemplate(w, "note-detail.html", note)
}

func main() {
	if err := ensureDataDir(); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/notes", getNotesHandler)
	http.HandleFunc("/note", createNoteHandler)
	http.HandleFunc("/note/{id}", getNoteHandler)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func init() {
	_ = io.Discard
}
