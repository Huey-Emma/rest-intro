package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type PostHandler struct {
	posts []Post
}

func (h PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid post id"))
		return
	}

	post := find(h.posts, func(p Post) bool { return p.ID == id })

	if post.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("post with id: %d not found", id)))
		return
	}

	b, err := json.Marshal(post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func find[TS ~[]T, T any](xs TS, fn func(T) bool) T {
	for _, i := range xs {
		if fn(i) {
			return i
		}
	}

	return *new(T)
}

func main() {
	posts := []Post{
		{ID: 1, Body: "this is post one"},
		{ID: 2, Body: "this is post two"},
	}

	router := chi.NewRouter()

	handler := PostHandler{posts}

	router.Get("/posts/{id}", handler.GetPost)

	server := &http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	_ = server.ListenAndServe()
}
