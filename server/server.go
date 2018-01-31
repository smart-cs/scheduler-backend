package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const logFormat = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Method}} {{.Path}}\n"

// Server runs the backend server.
type Server struct {
	Port            int
	Middleware      *negroni.Negroni
	ScheduleCreator ScheduleCreator
	AutoCompleter   AutoCompleter
}

// StandardResponse is the default response from the server.
type StandardResponse struct {
	OK     bool        `json:"OK"`
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// NewServer constructs a Server to listen on the given port.
func NewServer(port int) Server {
	server := Server{
		Port:            port,
		Middleware:      negroni.New(),
		ScheduleCreator: NewScheduleCreator(),
		AutoCompleter:   NewAutoCompleter(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/schedules", server.schedulesHandler).
		Methods("GET").
		Queries("courses", "{courses}")
	router.HandleFunc("/autocomplete", server.autocompleteHandler).
		Methods("GET").
		Queries("text", "{text}")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	logger := negroni.NewLogger()
	logger.SetDateFormat(time.Stamp)
	logger.SetFormat(logFormat)
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	server.Middleware.Use(logger)
	server.Middleware.Use(cors)
	server.Middleware.Use(negroni.NewRecovery())
	server.Middleware.UseHandler(router)
	return server
}

// Start starts the server.
func (s *Server) Start() {
	fmt.Printf("Listening on port %d\n", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Middleware)
}

func (s *Server) schedulesHandler(w http.ResponseWriter, r *http.Request) {
	courses := strings.Split(r.URL.Query().Get("courses"), ",")
	schedules := s.ScheduleCreator.Create(courses)
	s.respOK(w, schedules)
}

func (s *Server) autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	completes := s.AutoCompleter.CoursesWithPrefix(text)
	s.respOK(w, completes)
}

func (s *Server) respOK(w http.ResponseWriter, body interface{}) {
	r := StandardResponse{
		OK:     true,
		Status: http.StatusOK,
		Body:   body,
	}
	j, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic("can't marshal JSON")
	}
	w.Write(j)
}
