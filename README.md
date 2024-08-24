# Case Note Generator

This project is a web server that converts Markdown notes into Jira Wiki Notation. It utilizes the OpenAI API to generate text completions based on user prompts, serving static files and providing endpoints for generating completions and handling system prompts.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Frontend Functionality](#frontend-functionality)
- [License](#license)

## Features

- Convert Markdown notes to Jira Wiki Notation.
- Generate text completions using OpenAI's API.
- Serve static files (HTML, CSS, JavaScript).
- Handle file uploads for text (.txt) and Markdown (.md) files.
- Date selection for note generation.
- Copy generated content to clipboard.

## Requirements

- Go (version 1.22 or higher)
- OpenAI API key
- A `.env` file with the necessary environment variables

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/monstercameron/case-note-generator.git
   cd case-note-generator
   ```

2. Install dependencies:
   ```bash
   go get
   ```

3. Create a `.env` file in the root directory with the following content:
   ```plaintext
   OPENAI_API_KEY=your_openai_api_key
   OPENAI_API_MODEL=your_openai_model
   ```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```

2. Open your browser and navigate to `http://localhost:8080` to access the application.

3. Use the web interface to:
   - Upload text or Markdown files
   - Select a date
   - Generate Jira comments
   - Copy generated content to clipboard

## API Endpoints

- **GET /**: Serves the `index.html` file.
- **POST /generate**: Accepts a JSON body with a `prompt` field and returns a generated completion.
  - **Request Body**:
    ```json
    {
      "prompt": "Your prompt here"
    }
    ```
  - **Response**:
    ```json
    {
      "completion": "Generated text here"
    }
    ```

- **GET /systemprompt**: Returns a message indicating the GET endpoint was hit.
- **POST /systemprompt**: Returns a message indicating the POST endpoint was hit.

## Frontend Functionality

The application includes a JavaScript file (`static/script/script.js`) that provides the following features:

- Drag and drop file upload
- File selection via input
- Date selection
- Jira comment generation
- Copying generated content to clipboard

## License

This project is licensed under the MIT License. See the LICENSE file for details.