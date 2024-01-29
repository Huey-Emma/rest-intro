package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Cache struct {
	data map[string]interface{}
	*sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value

	if ttl > 0 {
		time.AfterFunc(ttl, func() {
			c.Del(key)
		})
	}
}

func (c *Cache) Del(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
}

type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Rating int    `json:"rating"`
}

type Store struct {
	posts []Post
}

func (s *Store) GetPosts() []Post {
	time.Sleep(3 * time.Second)
	return s.posts
}

func (s *Store) GetPopularPosts() []Post {
	time.Sleep(3 * time.Second)
	return filter(s.posts, func(p Post) bool { return p.Rating >= 4 })
}

type Storer interface {
	GetPosts() []Post
	GetPopularPosts() []Post
}

type PostHandler struct {
	store Storer
	cache *Cache
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts := h.store.GetPosts()
	b, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h *PostHandler) GetPopularPosts(w http.ResponseWriter, r *http.Request) {
	posts := h.store.GetPopularPosts()
	key := strings.TrimPrefix(r.URL.Path, "/")
	h.cache.Set(key, posts, 10*time.Second)

	b, _ := json.Marshal(posts)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Middlewares
func JSONContentType(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		n.ServeHTTP(w, r)
	})
}

func CacheMiddleware(cache *Cache) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := strings.TrimPrefix(r.URL.Path, "/")

			if cached, ok := cache.Get(key); ok {
				b, _ := json.Marshal(cached)
				w.WriteHeader(http.StatusOK)
				w.Write(b)
				return
			}

			n.ServeHTTP(w, r)
		})
	}
}

// Helpers
func filter[TS ~[]T, T interface{}](xs TS, fn func(T) bool) TS {
	out := make(TS, 0)

	for _, i := range xs {
		if fn(i) {
			out = append(out, i)
		}
	}

	return out
}

func main() {
	cache := &Cache{
		data:    make(map[string]interface{}),
		RWMutex: &sync.RWMutex{},
	}

	store := &Store{
		posts: []Post{
			{Title: "Post one", Body: "Post one content", Rating: 4},
			{Title: "Post two", Body: "Post two content", Rating: 2},
			{Title: "Post three", Body: "Post three content", Rating: 5},
			{Title: "Post four", Body: "Post four content", Rating: 3},
		},
	}

	router := chi.NewRouter()
	router.Use(JSONContentType)

	handler := &PostHandler{store, cache}
	router.Get("/posts", handler.GetPosts)
	router.With(CacheMiddleware(cache)).Get("/posts/most-popular", handler.GetPopularPosts)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	_ = server.ListenAndServe()
}
