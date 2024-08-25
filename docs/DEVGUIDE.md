// ... existing content ...

## Developer Guide

### How main.go Works

The `main.go` file is the entry point of the application. Here's an overview of its structure, functionality, and key third-party imports:

1. **Imports**:
   - Standard library imports for HTTP handling, JSON processing, logging, etc.
   - Third-party packages:
     - `github.com/joho/godotenv`: Used to load environment variables from a .env file.
     - `github.com/sashabaranov/go-openai`: Provides a client for interacting with the OpenAI API.

2. **Types**: 
   - Defines custom types like `GenerateRequest` for handling API requests.

3. **Middleware**:
   - `logMiddleware`: Logs the start and completion of each request.

4. **Main function**:
   - Uses `godotenv.Load()` to load environment variables from the .env file.
   - Initializes the OpenAI client using `openai.NewClient()` with the API key from environment variables.
   - Sets up HTTP routes and starts the server with graceful shutdown capabilities.

5. **Route Handlers**:
   - `indexHandler`: Serves the main HTML page.
   - `generateHandler`: Handles the generation of completions using the OpenAI client.
   - `systemPromptHandler`: Handles GET requests for the system prompt.
   - `systemPromptPostHandler`: Handles POST requests to update the system prompt.
   - `healthHandler`: Provides a health check endpoint.

6. **Helper Functions**:
   - `logSetup`: Configures logging to both console and file.
   - `getCompletion`: Uses the OpenAI client to retrieve completions from the API.

### Third-Party Package Usage

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

### Environment Variables

The application requires a `.env` file in the root directory with the following variables:

```
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_MODEL=gpt-3.5-turbo
PORT=8080
```

Make sure to replace `your_openai_api_key_here` with your actual OpenAI API key.

### Frontend

The frontend code is located in the `static` folder. The main HTML file is `index.html`, which renders the user interface for interacting with the AI model. The `script.js` file contains the JavaScript code that handles user input, sends requests to the backend, and updates the UI with the AI-generated responses.

### How to Add More Endpoints

To extend the project with additional functionality, you can add new endpoints. Here's a step-by-step guide:

1. **Define a new handler function** in `main.go`:
   ```go
   func newEndpointHandler(w http.ResponseWriter, r *http.Request) {
       // Implement your new endpoint logic here
   }
   ```

2. **Add the new route** in the `main()` function:
   ```go
   http.HandleFunc("GET/api/new-endpoint", newEndpointHandler)
   ```

3. **Implement the handler logic**, considering:
   - Request parsing (e.g., JSON decoding for POST requests)
   - Error handling
   - Response formatting (e.g., JSON encoding)
   - Logging

4. **Update the frontend** in `static/script.js` to interact with the new endpoint:
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

6. **Document the new endpoint** in this README, including its purpose, request/response format, and any new environment variables or dependencies.

Remember to consider security implications, rate limiting, and performance impact when adding new endpoints.

Remember to handle errors, add appropriate logging, and update tests when extending the application. Also, consider the impact on performance and security when implementing new features.