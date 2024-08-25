// ... existing content ...

## Developer Guide

### How main.go Works

The `main.go` file is the entry point of the application. Here's an overview of its structure, functionality, and key third-party imports:

1. **Imports**:
   - Standard library imports for HTTP handling, JSON processing, logging, etc.
   - Third-party packages:
     - `github.com/joho/godotenv`: Used to load environment variables from a .env file.
     - `github.com/sashabaranov/go-openai`: Provides a client for interacting with the OpenAI API.
     - `github.com/shirou/gopsutil/cpu`: Used for CPU usage information in the health check.

2. **Types**: 
   - Defines custom types like `GenerateRequest`, `SummaryRequest`, and `HealthData` for handling various API requests and responses.

3. **Middleware**:
   - `logMiddleware`: Logs the start and completion of each request.

4. **Main function**:
   - Uses `godotenv.Load()` to load environment variables from the .env file.
   - Initializes the OpenAI client using `openai.NewClient()` with the API key from environment variables.
   - Reads system prompts for notes and summaries from files.
   - Sets up HTTP routes and starts the server with graceful shutdown capabilities.

5. **Route Handlers**:
   - `indexHandler`: Serves the main HTML page.
   - `generateHandler`: Handles the generation of Jira comments using the OpenAI client.
   - `summaryHandler`: Handles the generation of summaries using the OpenAI client.
   - `systemPromptHandler`: Handles GET requests for the system prompt.
   - `systemPromptPostHandler`: Handles POST requests to update the system prompt.
   - `healthHandler`: Provides a health check endpoint with detailed system information.
   - `logFileHandler`: Serves the contents of the log file.

6. **Helper Functions**:
   - `logSetup`: Configures logging to both console and file.
   - `getCompletion`: Uses the OpenAI client to retrieve completions from the API.

## Third-Party Package Usage

1. **godotenv**:
   - Used in the `main()` function to load environment variables:
     ```go
     err := godotenv.Load()
     if err != nil {
         log.Fatal("Error loading .env file:", err)
     }
     ```

2. **go-openai**:
   - Used to create an OpenAI client and make API requests:
     ```go
     client := openai.NewClient(openaiAPIKey)
     ```
   - The `getCompletion()` function uses this client to create chat completions:
     ```go
     resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{...})
     ```

3. **gopsutil**:
   - Used in the `healthHandler` to get CPU usage information:
     ```go
     cpuUsage, err := cpu.Percent(0, false)
     ```

## Environment Variables

The application requires a `.env` file in the root directory with the following variables:

```
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_MODEL=gpt-3.5-turbo
PORT=8080
```

Make sure to replace `your_openai_api_key_here` with your actual OpenAI API key.

## Frontend

The frontend code is located in the `static` folder. The main HTML file is `index.html`, which renders the user interface for interacting with the AI model. The `script.js` file contains the JavaScript code that handles user input, sends requests to the backend, and updates the UI with the AI-generated responses.

Key frontend features include:
- Drag and drop file upload
- Date selection for note generation
- Jira comment and summary generation
- Copying generated content to clipboard
- Editing and updating system prompts
- Viewing and refreshing application logs
- Collapsible sections for better organization

## API Collection

Here's a detailed collection of all the API endpoints in the application:

### 1. Serve Index Page
- **Method**: GET
- **Path**: /
- **Handler**: `indexHandler`
- **Description**: Serves the main index.html file.
- **Response**: HTML content

### 2. Generate Jira Comment
- **Method**: POST
- **Path**: /generate
- **Handler**: `generateHandler`
- **Description**: Generates Jira comments based on uploaded files and selected date.
- **Request Body**:
  ```json
  {
    "prompt": "string",
    "date": "string"
  }
  ```
- **Response**:
  ```json
  {
    "completion": "string"
  }
  ```

### 3. Generate Summary
- **Method**: POST
- **Path**: /summary
- **Handler**: `summaryHandler`
- **Description**: Generates summaries based on uploaded files.
- **Request Body**:
  ```json
  {
    "prompt": "string"
  }
  ```
- **Response**:
  ```json
  {
    "summary": "string"
  }
  ```

### 4. Get System Prompt
- **Method**: GET
- **Path**: /systemprompt
- **Handler**: `systemPromptHandler`
- **Description**: Retrieves system prompts or specific prompt content.
- **Query Parameters**:
  - `file`: (optional) Name of the specific prompt file to retrieve
- **Response**:
  - If `file` parameter is provided:
    ```json
    {
      "file": "string",
      "content": "string"
    }
    ```
  - If `file` parameter is not provided:
    ```json
    ["prompt1.txt", "prompt2.txt", ...]
    ```

### 5. Update System Prompt
- **Method**: POST
- **Path**: /systemprompt
- **Handler**: `systemPromptPostHandler`
- **Description**: Updates a system prompt.
- **Request Body**:
  ```json
  {
    "filename": "string",
    "prompt": "string"
  }
  ```
- **Response**:
  ```json
  {
    "status": "string",
    "message": "string"
  }
  ```

### 6. Health Check
- **Method**: GET
- **Path**: /health
- **Handler**: `healthHandler`
- **Description**: Returns health check data about the application.
- **Response**:
  ```json
  {
    "status": "string",
    "version": "string",
    "uptime": "string",
    "goVersion": "string",
    "numGoroutine": "number",
    "cpuUsage": "number",
    "cpuCount": "number",
    "memUsageMB": "number"
  }
  ```

### 7. Get Logs
- **Method**: GET
- **Path**: /logs
- **Handler**: `logFileHandler`
- **Description**: Retrieves application logs.
- **Response**: Plain text content of the log file

### 8. Serve Static Files
- **Method**: GET
- **Path**: /static/
- **Handler**: http.FileServer
- **Description**: Serves static files (HTML, CSS, JavaScript) from the `static` directory.

Remember to include appropriate error handling and authentication/authorization mechanisms where necessary when implementing or using these endpoints.

## How to Add More Endpoints

To extend the project with additional functionality, you can add new endpoints. Here's a step-by-step guide:

1. **Define a new handler function** in `main.go`:
   ```go
   func newEndpointHandler(w http.ResponseWriter, r *http.Request) {
       // Implement your new endpoint logic here
   }
   ```

2. **Add the new route** in the `main()` function:
   ```go
   http.HandleFunc("GET /api/new-endpoint", logMiddleware(newEndpointHandler))
   ```

3. **Implement the handler logic**, considering:
   - Request parsing (e.g., JSON decoding for POST requests)
   - Error handling
   - Response formatting (e.g., JSON encoding)
   - Logging

4. **Update the frontend** in `static/script/script.js` to interact with the new endpoint:
   ```javascript
   async function callNewEndpoint() {
       const response = await fetch('/api/new-endpoint', {
           method: 'POST',
           headers: { 'Content-Type': 'application/json' },
           body: JSON.stringify({ /* request data */ }),
       });
       const data = await response.json();
       // Handle the response
   }
   ```

5. **Add appropriate error handling and logging** throughout the new code.

6. **Document the new endpoint** in the README.md, including its purpose, request/response format, and any new environment variables or dependencies.

Remember to consider security implications, rate limiting, and performance impact when adding new endpoints. Also, handle errors, add appropriate logging, and update tests when extending the application. Consider the impact on performance and security when implementing new features.

## Logging

The application uses a custom logging setup that writes logs to both the console and a file (`logs/log.log`). The `logSetup()` function in `main.go` configures this logging system.

## Health Check

The `healthHandler` provides detailed information about the application's status, including:
- Application status
- Version
- Uptime
- Go version
- Number of goroutines
- CPU usage
- CPU count
- Memory usage

This endpoint can be useful for monitoring the application's performance and health.

## System Prompts

The application uses separate system prompts for generating Jira comments and summaries. These prompts are stored in files (`static/document/notes.prompt` and `static/document/summary.prompt`) and can be updated through the frontend interface.

Remember to thoroughly test any new features or endpoints before deploying them to production.