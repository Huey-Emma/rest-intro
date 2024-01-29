package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type Post struct {
	ID    int64    `json:"id"`
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].Title < p[j].Title }

type PostHandler struct {
	posts Posts
}

func (h PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts Posts

	rawQ, _ := url.ParseQuery(r.URL.RawQuery)

	sortQ := rawQ.Get("sort")
	tagQ := rawQ.Get("tag")
	offset, _ := strconv.Atoi(rawQ.Get("page"))
	limit, _ := strconv.Atoi(rawQ.Get("pageSize"))

	if tagQ != "" {
		posts = filter(
			h.posts,
			func(t Post) bool {
				return some(t.Tags, func(t string) bool {
					return t == tagQ
				})
			},
		)
	} else {
		posts = h.posts
	}

	if sortQ == "title" {
		sort.Sort(posts)
	} else if sortQ == "-title" {
		sort.Sort(sort.Reverse(posts))
	}

	var pageNum, pageSize int

	if limit == 0 {
		pageSize = 5
	} else {
		pageSize = limit
	}

	if offset > 0 {
		pageNum = offset - 1
	}

	start, end := paginate(pageNum, pageSize, len(posts))
	b, _ := json.Marshal(posts[start:end])
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Middlewares
func JSONContentType(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		n.ServeHTTP(w, r)
	})
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

func some[TS ~[]T, T comparable](xs TS, fn func(T) bool) bool {
	for _, i := range xs {
		if fn(i) {
			return true
		}
	}
	return false
}

func paginate(pageNum, pageSize, lenght int) (int, int) {
	var start, end int

	start = pageNum * pageSize

	if start > lenght {
		start = lenght
	}

	end = start + pageSize

	if end > lenght {
		end = lenght
	}

	return start, end
}

func randint(mini, maxi int) int {
	return rand.Intn((maxi-mini)+1) + mini
}

func alphabets() string {
	builder := strings.Builder{}

	for i := 'a'; i < 'a'+26; i++ {
		builder.WriteRune(i)
	}

	for i := 'A'; i < 'A'+26; i++ {
		builder.WriteRune(i)
	}

	return builder.String()
}

func randomString(l int) string {
	builder := strings.Builder{}
	opts := alphabets()

	for i := 0; i < l; i++ {
		builder.WriteByte(opts[rand.Intn(len(opts)-1)])
	}

	return builder.String()
}

func main() {
	router := chi.NewRouter()

	posts := make(Posts, 100)

	for i := range posts {
		p := Post{
			ID:    time.Now().UTC().UnixNano(),
			Title: randomString(10),
			Body:  randomString(100),
		}

		if i%2 == 0 {
			p.Tags = append(p.Tags, []string{"fiction", "adventure"}...)
		} else {
			p.Tags = append(p.Tags, []string{"autobiography", "motivation"}...)
		}

		posts[i] = p
	}

	handler := PostHandler{posts}

	router.Use(JSONContentType)
	router.Get("/posts", handler.GetPosts)

	server := &http.Server{
		Addr:    ":6000",
		Handler: router,
	}

	_ = server.ListenAndServe()
}
