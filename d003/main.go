package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetComments(w http.ResponseWriter, _ *http.Request) {
	comments := []map[string]interface{}{
		{"body": "this was an awesome read"},
		{"body": "this was a thoughtless piece, the author is a dolt"},
		{"body": "I am just here for the violence"},
	}

	b, _ := json.Marshal(&comments)
	w.Write(b)
}

func main() {
	// initialize router
	router := chi.NewRouter()

	// gets all comments associted a post with the given id
	router.Get("/posts/{id}/comments", GetComments)

	// configure server
	server := &http.Server{
		Addr:    ":9999",
		Handler: router,
	}

	_ = server.ListenAndServe()
}
