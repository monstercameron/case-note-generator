package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/joho/godotenv"
)

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	client := openai.NewClient("your token")

	// Get the OpenAI API key and model from environment variables
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	openaiAPIModel := os.Getenv("OPENAI_API_MODEL")
	if openaiAPIKey == "" || openaiAPIModel == "" {
		log.Fatal("OPENAI_API_KEY or OPENAI_API_MODEL is not set in the environment variables")
	}

	// handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes using the new style
	http.HandleFunc("GET /", indexHandler) // GET /
	http.HandleFunc("POST /generate", func(w http.ResponseWriter, r *http.Request) {
		generateHandler(w, r, client, openaiAPIModel) // POST /generate
	})
	http.HandleFunc("GET /systemprompt", systemPromptHandler)      // GET /systemprompt
	http.HandleFunc("POST /systemprompt", systemPromptPostHandler) // POST /systemprompt

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

// generateHandler handles the generation of completions
func generateHandler(w http.ResponseWriter, r *http.Request, client *openai.Client, model string) {
	var req GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get completion from OpenAI
	completion, err := getCompletion(client, model, req.Prompt)
	if err != nil {
		http.Error(w, "Error getting completion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the completion as JSON
	response := map[string]string{"completion": completion}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getCompletion retrieves a completion from the OpenAI API
func getCompletion(client *openai.Client, model string, prompt string) (string, error) {
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	// Extract and return the completion text
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no completion found")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html") // Serve index.html
}

func systemPromptHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET logic here
	fmt.Fprintln(w, "System prompt GET endpoint hit")
}

func systemPromptPostHandler(w http.ResponseWriter, r *http.Request) {
	// Handle POST logic here
	fmt.Fprintln(w, "System prompt POST endpoint hit")
}
