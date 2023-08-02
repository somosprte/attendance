package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("frontend/build/static"))))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "frontend/build/index.html") })

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/meetings", createMeeting)
		r.Post("/register/{meetingID}", registerAttendee) // Trocado de "registerParticipant" para "registerAttendee"
		r.Get("/meetings/{meetingID}", getMeeting)
	})

	fmt.Println("Server running on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
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
