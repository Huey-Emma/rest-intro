package main

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	version, _ := r.Context().Value("version").(int)

	switch version {
	case 1:
		h.HelloWorldV1(w, r)
		return
	case 2:
		h.HelloWorldV2(w, r)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"detail": "invalid version"}`))
	}
}

func (h *HelloHandler) HelloWorldV1(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world\"}"))
}

func (h *HelloHandler) HelloWorldV2(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Hello world 😏\"}"))
}

func VersionMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var version int

		p := regexp.MustCompile(`/api/(v\d+)/.*`)
		q, _ := url.ParseQuery(r.URL.RawQuery)

		if version == 0 && p.MatchString(r.URL.Path) {
			m := p.FindStringSubmatch(r.URL.Path)
			version, _ = strconv.Atoi(strings.TrimPrefix(m[1], "v"))
		}

		if version == 0 && q != nil {
			version, _ = strconv.Atoi(q.Get("version"))
		}

		if version == 0 {
			version, _ = strconv.Atoi(strings.TrimPrefix(r.Header.Get("Accept-version"), "v"))
		}

		ctx := context.WithValue(r.Context(), "version", version)

		n.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JSONContentType(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		n.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(VersionMiddleware, JSONContentType)

	r.Get("/api*", (&HelloHandler{}).ServeHTTP)

	server := http.Server{
		Addr:    ":9999",
		Handler: r,
	}

	_ = server.ListenAndServe()
}
