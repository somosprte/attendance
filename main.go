package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Meeting struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Time        string   `json:"time"`
	Description string   `json:"description"`
	Attendees   []string `json:"attendees"` // Trocado de "participants" para "attendees"
}

var meetings = make(map[string]*Meeting)
var meetingsMutex = &sync.Mutex{}

//go:embed frontend/build/*
var content embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fsys, _ := fs.Sub(content, "frontend/build")
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(fsys))))

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/meetings", createMeeting)
		r.Get("/meetings", getMeetings)
		r.Post("/register/{meetingID}", registerAttendee) // Trocado de "registerParticipant" para "registerAttendee"
		r.Get("/meetings/{meetingID}", getMeeting)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Porta padrão se não for encontrada porta no ENV
	}

	fmt.Println("Server running on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func createMeeting(w http.ResponseWriter, r *http.Request) {
	var meeting Meeting
	if err := json.NewDecoder(r.Body).Decode(&meeting); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	meeting.ID = uuid.New().String()
	meeting.Attendees = []string{} // Trocado de "Participants" para "Attendees"

	meetingsMutex.Lock()
	meetings[meeting.ID] = &meeting
	meetingsMutex.Unlock()

	json.NewEncoder(w).Encode(meeting)
}

func registerAttendee(w http.ResponseWriter, r *http.Request) { // Trocado de "registerParticipant" para "registerAttendee"
	meetingID := chi.URLParam(r, "meetingID")
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	meetingsMutex.Lock()
	meeting, exists := meetings[meetingID]
	if !exists {
		meetingsMutex.Unlock()
		http.Error(w, "Meeting not found", http.StatusNotFound)
		return
	}
	meeting.Attendees = append(meeting.Attendees, name) // Trocado de "Participants" para "Attendees"
	meetingsMutex.Unlock()

	json.NewEncoder(w).Encode(meeting)
}

func getMeeting(w http.ResponseWriter, r *http.Request) {
	meetingID := chi.URLParam(r, "meetingID")

	meetingsMutex.Lock()
	meeting, exists := meetings[meetingID]
	meetingsMutex.Unlock()
	if !exists {
		http.Error(w, "Meeting not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(meeting)
}

func getMeetings(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(meetings)
}
