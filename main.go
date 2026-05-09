package main

import (
	"log"

	"generatePDF/cores"

	"github.com/valyala/fasthttp"
)

func main() {
	r := cores.NewRouter()

	port := ":8080"
	log.Printf("Starting server on %s", port)

	if err := fasthttp.ListenAndServe(port, r.Handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
