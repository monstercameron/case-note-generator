package main

import (
    "fmt"
    "net/http"
    "github.com/sashabaranov/go-openai"
)

func main() {
    http.HandleFunc("GET /", indexHandler)                     // GET /
    http.HandleFunc("POST /generate", generateHandler)          // POST /generate
    http.HandleFunc("GET /systemprompt", systemPromptHandler)    // GET /systemprompt
    http.HandleFunc("POST /systemprompt", systemPromptPostHandler) // POST /systemprompt

    fmt.Println("Server is running on :8080")
    http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html") // Serve index.html
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
    // Handle generation logic here
    fmt.Fprintln(w, "Generate endpoint hit")
}

func systemPromptHandler(w http.ResponseWriter, r *http.Request) {
    // Handle GET logic here
    fmt.Fprintln(w, "System prompt GET endpoint hit")
}

func systemPromptPostHandler(w http.ResponseWriter, r *http.Request) {
    // Handle POST logic here
    fmt.Fprintln(w, "System prompt POST endpoint hit")
}