package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Date   string `json:"date"`
}

// logMiddleware creates a middleware that logs all API calls
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the OpenAI API key and model from environment variables
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	openaiAPIModel := os.Getenv("OPENAI_API_MODEL")
	if openaiAPIKey == "" || openaiAPIModel == "" {
		log.Fatal("OPENAI_API_KEY or OPENAI_API_MODEL is not set in the environment variables")
	}
	client := openai.NewClient(openaiAPIKey)

	// Read system prompt from file
	systemPrompt, err := os.ReadFile("static/document/system.prompt")
	if err != nil {
		log.Fatal("Error reading system prompt file:", err)
	}

	// handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Define routes using the new style with logging middleware
	http.HandleFunc("GET /", logMiddleware(indexHandler))
	http.HandleFunc("POST /generate", logMiddleware(func(w http.ResponseWriter, r *http.Request) {
		generateHandler(w, r, client, openaiAPIModel, string(systemPrompt))
	}))
	http.HandleFunc("GET /systemprompt", logMiddleware(systemPromptHandler))
	http.HandleFunc("POST /systemprompt", logMiddleware(systemPromptPostHandler))

	// Get the PORT from environment variables, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil, // Use the default ServeMux
	}

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server is running on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Set up channel to listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// generateHandler handles the generation of completions
func generateHandler(w http.ResponseWriter, r *http.Request, client *openai.Client, model string, systemPrompt string) {
	var req GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	prePrompt := fmt.Sprintf("Create Case notes for %s. Generate a Jira comment based on the following information:", req.Date)
	prompt := prePrompt + "\n\n" + req.Prompt

	// Get completion from OpenAI
	completion, err := getCompletion(client, model, prompt, systemPrompt)
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
func getCompletion(client *openai.Client, model string, prompt string, systemPrompt string) (string, error) {
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func systemPromptHandler(w http.ResponseWriter, r *http.Request) {
	// Read the system prompt file
	systemPrompt, err := os.ReadFile("static/document/system.prompt")
	if err != nil {
		http.Error(w, "Error reading system prompt file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response struct
	response := struct {
		SystemPrompt string `json:"systemPrompt"`
	}{
		SystemPrompt: string(systemPrompt),
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func systemPromptPostHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to parse the incoming JSON
	var requestBody struct {
		Prompt string `json:"prompt"`
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error parsing JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Write the new prompt to the file
	err = os.WriteFile("static/document/system.prompt", []byte(requestBody.Prompt), 0644)
	if err != nil {
		http.Error(w, "Error writing to system prompt file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "System prompt updated successfully",
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}