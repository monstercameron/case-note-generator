package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/shirou/gopsutil/cpu"
)

// GenerateRequest represents the structure of the incoming JSON request for generation.
type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Date   string `json:"date"`
}

// logMiddleware creates a middleware that logs all API calls.
// It wraps the given http.HandlerFunc and logs the start and completion of each request.
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
	}
}

// main is the entry point of the application.
// It sets up logging, loads environment variables, initializes the OpenAI client,
// sets up HTTP routes, and starts the server with graceful shutdown capabilities.
func main() {
	// Set up logging to both console and file
	logSetup()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	log.Println("Environment variables loaded successfully")

	// Get the OpenAI API key and model from environment variables
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	openaiAPIModel := os.Getenv("OPENAI_API_MODEL")
	if openaiAPIKey == "" || openaiAPIModel == "" {
		log.Fatal("OPENAI_API_KEY or OPENAI_API_MODEL is not set in the environment variables")
	}
	client := openai.NewClient(openaiAPIKey)
	log.Println("OpenAI client initialized")

	// Read system prompt from file
	systemPrompt, err := os.ReadFile("static/document/system.prompt")
	if err != nil {
		log.Fatal("Error reading system prompt file:", err)
	}
	log.Println("System prompt loaded successfully")

	// Handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Define routes using the new style with logging middleware
	http.HandleFunc("GET /", logMiddleware(indexHandler))
	http.HandleFunc("POST /generate", logMiddleware(func(w http.ResponseWriter, r *http.Request) {
		generateHandler(w, r, client, openaiAPIModel, string(systemPrompt))
	}))
	http.HandleFunc("GET /systemprompt", logMiddleware(systemPromptHandler))
	http.HandleFunc("POST /systemprompt", logMiddleware(systemPromptPostHandler))

	// Add the health check endpoint
	http.HandleFunc("GET /health", logMiddleware(healthHandler))

	// Get the PORT from environment variables, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT not set, defaulting to 8080")
	}

	// Create a new server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil, // Use the default ServeMux
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server is running on :%s\n", port)
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

// logSetup configures logging to both console and file.
// It creates a logs directory if it doesn't exist and sets up a multi-writer
// to output logs to both stdout and a log file.
func logSetup() {
	// Check if /logs directory exists, create if not
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal("Error creating logs directory:", err)
		}
		log.Println("Logs directory created")
	}

	// Open log file (create if not exists)
	logFile, err := os.OpenFile("logs/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	// Create multi writer for both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set log output to multi writer
	log.SetOutput(multiWriter)

	// Set log flags for more detailed logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Logging setup completed")
}

// generateHandler handles the generation of completions.
// It decodes the incoming JSON request, prepares the prompt,
// calls the OpenAI API for completion, and returns the result as JSON.
func generateHandler(w http.ResponseWriter, r *http.Request, client *openai.Client, model string, systemPrompt string) {
	var req GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received generate request for date: %s", req.Date)

	// clear console
	// fmt.Println("\033[H\033[2J")
	prePrompt := fmt.Sprintf("Create a Jira Comment based on the following information strictly and only for the date of %s:", req.Date)
	prompt := prePrompt + "\n" + req.Prompt

	// print prompt for debugging
	// fmt.Println(prompt)

	// Get completion from OpenAI
	completion, err := getCompletion(client, model, prompt, systemPrompt)
	if err != nil {
		log.Printf("Error getting completion: %v", err)
		http.Error(w, "Error getting completion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Completion generated successfully")

	// Return the completion as JSON
	response := map[string]string{"completion": completion}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("Generate response sent successfully")
}

// getCompletion retrieves a completion from the OpenAI API.
// It sends a request to the OpenAI API with the given prompt and system prompt,
// and returns the generated completion text.
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
		log.Printf("Error creating chat completion: %v", err)
		return "", err
	}

	// Extract and return the completion text
	if len(resp.Choices) > 0 {
		log.Println("Chat completion received successfully")
		return resp.Choices[0].Message.Content, nil
	}
	log.Println("No completion found in the response")
	return "", fmt.Errorf("no completion found")
}

// indexHandler serves the main index.html file.
// It checks if the requested path is "/" and serves the index.html file accordingly.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("404 Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	log.Println("Serving index.html")
	http.ServeFile(w, r, "static/index.html")
}

// systemPromptHandler handles GET requests for the system prompt.
// It reads the system prompt from a file and returns it as a JSON response.
func systemPromptHandler(w http.ResponseWriter, r *http.Request) {
	// Read the system prompt file
	systemPrompt, err := os.ReadFile("static/document/system.prompt")
	if err != nil {
		log.Printf("Error reading system prompt file: %v", err)
		http.Error(w, "Error reading system prompt file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("System prompt read successfully")

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
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("System prompt sent successfully")
}

// systemPromptPostHandler handles POST requests to update the system prompt.
// It receives a new prompt in the request body, writes it to the system prompt file,
// and returns a success message as a JSON response.
func systemPromptPostHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to parse the incoming JSON
	var requestBody struct {
		Prompt string `json:"prompt"`
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error parsing JSON request: %v", err)
		http.Error(w, "Error parsing JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Received new system prompt")

	// Write the new prompt to the file
	err = os.WriteFile("static/document/system.prompt", []byte(requestBody.Prompt), 0644)
	if err != nil {
		log.Printf("Error writing to system prompt file: %v", err)
		http.Error(w, "Error writing to system prompt file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("System prompt updated successfully")

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
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("System prompt update response sent successfully")
}

// HealthData represents the structure of the health check response.
type HealthData struct {
	Status       string  `json:"status"`
	Version      string  `json:"version"`
	Uptime       string  `json:"uptime"`
	GoVersion    string  `json:"goVersion"`
	NumGoroutine int     `json:"numGoroutine"`
	CpuUsage     float64 `json:"cpuUsage"`
	CpuCount     int     `json:"cpuCount"`
	MemUsageMB   float64 `json:"memUsageMB"`
}

var startTime = time.Now()

// healthHandler handles GET requests for the health check endpoint.
// It returns a JSON payload with relevant health data about the application.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get CPU usage
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		cpuUsage = []float64{0} // Default to 0 if there's an error
	}

	health := HealthData{
		Status:       "OK",
		Version:      "1.0.0", // You should replace this with your actual version
		Uptime:       time.Since(startTime).String(),
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		CpuCount:     runtime.NumCPU(),
		CpuUsage:     cpuUsage[0],
		MemUsageMB:   float64(m.Alloc) / 1024 / 1024, // Convert bytes to MB
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("Error encoding health check JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("Health check response sent successfully")
}