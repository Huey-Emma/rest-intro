package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	PRODUCTS_REGEX = regexp.MustCompile(`^/products$`)
	PRODUCT_REGEX  = regexp.MustCompile(`^/products/([0-9]+)$`)
)

type Product struct {
	ID    int64
	Name  string
	Price float64
}

func (p Product) String() string {
	return fmt.Sprintf("Product id=%d name=%s price=%.2f", p.ID, p.Name, p.Price)
}

type ProductHandler struct {
	products []Product
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && PRODUCTS_REGEX.MatchString(r.URL.Path):
		h.GetProducts(w, r)
		return
	case r.Method == http.MethodPost && PRODUCTS_REGEX.MatchString(r.URL.Path):
		h.CreateProduct(w, r)
		return
	case r.Method == http.MethodPut && PRODUCT_REGEX.MatchString(r.URL.Path):
		h.UpdateProduct(w, r)
		return
	case r.Method == http.MethodDelete && PRODUCT_REGEX.MatchString(r.URL.Path):
		h.DeleteProduct(w, r)
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(h.products)
	w.Write(b)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var p Product
	_ = json.Unmarshal(b, &p)
	p.ID = time.Now().UTC().UnixNano()

	h.products = append(h.products, p)

	b, _ = json.Marshal(p)
	w.Write(b)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updated product"))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Set JSON Content-Type middleware
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
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
	products := make([]Product, 10)

	for idx := range products {
		products[idx] = Product{
			ID:    time.Now().UTC().UnixNano(),
			Name:  randomString(randint(5, 15)),
			Price: float64(randint(500, 10000)),
		}
	}

	productHandler := JSONContentType(&ProductHandler{products})

	http.Handle("/products", productHandler)
	http.Handle("/products/", productHandler)

	http.ListenAndServe(":9000", nil)
}
