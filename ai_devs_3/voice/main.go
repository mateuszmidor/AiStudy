package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func uploadAudioHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit your max input length!
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Save the audio data to a file (you can customize the filename)
	err = ioutil.WriteFile("uploads/recording.wav", data, 0644)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", "recording.wav")
}

func main() {
	// Create uploads directory if it doesn't exist
	os.MkdirAll("uploads", os.ModePerm)

	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadAudioHandler).Methods("POST")

	// Allow CORS from anywhere
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("Server starting at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
