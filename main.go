package main

import (
    "log"
    "net/http"
    "github.com/dgraph-io/badger/v4"
    "github.com/RobAntunes/tig/internal/api"
    "github.com/RobAntunes/tig/internal/intent/storage"
)

func main() {
    // Initialize BadgerDB
    db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize repositories
    intentStore := storage.NewBadgerStore(db)
    intentHandler := api.NewIntentHandler(intentStore)

    // Set up router
    mux := http.NewServeMux()

    // Health checks
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Tig is running!"))
    })

    // Intent endpoints
    mux.HandleFunc("POST /api/intents", intentHandler.Create)
    mux.HandleFunc("GET /api/intents/{id}", intentHandler.Get)
    mux.HandleFunc("PUT /api/intents/{id}", intentHandler.Update)

    log.Printf("Starting Tig on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}