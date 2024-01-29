package main

import (
	"log"
	"net/http"
)

// Handler function
func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message:\": \"hello world\"}"))
}

// Set JSON Content-Type middleware
func JSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/", JSONContentType(HelloHandler))

	log.Println("server is listening")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

// Express JS Equivalent
// import express from 'express'
//
// app = express()
//
// app.use(express.json())
//
// app.get('/', (res, req) => {
//    res.json({message: 'hello world'})
// })
//
// app.listen(8000, () => {
//		console.log("server is listening")
// })
