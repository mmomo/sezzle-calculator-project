package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mmomo/sezzle-calculator-project/backend/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/calculate", api.CalculateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
