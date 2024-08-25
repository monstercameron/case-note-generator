# Case Note Generator

This project is a web server that converts Markdown notes into Jira Wiki Notation and generates summaries. It utilizes the OpenAI API to generate text completions based on user prompts, serving static files and providing endpoints for generating completions, summaries, and handling system prompts.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Frontend Functionality](#frontend-functionality)
- [Developer Guide](#developer-guide)
- [License](#license)

## Features

- Convert Markdown notes to Jira Wiki Notation.
- Generate text completions and summaries using OpenAI's API.
- Serve static files (HTML, CSS, JavaScript).
- Handle file uploads for text (.txt) and Markdown (.md) files.
- Date selection for note generation.
- Copy generated content to clipboard.
- Edit and update system prompts.
- View and refresh application logs.
- Health check endpoint for monitoring application status.

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
   OPENAI_API_MODEL=gpt-3.5-turbo
   PORT=8080 # Optional, defaults to 8080 if not set
   ```

   Replace `your_openai_api_key` with your actual OpenAI API key.

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```

2. Open your browser and navigate to `http://localhost:8080` to access the application.

3. Use the web interface to:
   - Upload text or Markdown files
   - Select a date for note generation
   - Generate Jira comments and summaries
   - Copy generated content to clipboard
   - Edit and update system prompts
   - View and refresh application logs

## API Endpoints

- **GET /**: Serves the `index.html` file.
- **POST /generate**: Generates Jira comments based on uploaded files and selected date.
- **POST /summary**: Generates summaries based on uploaded files.
- **GET /systemprompt**: Retrieves system prompts or specific prompt content.
- **POST /systemprompt**: Updates a system prompt.
- **GET /health**: Returns health check data about the application.
- **GET /logs**: Retrieves application logs.

For more detailed information about the API endpoints and how to add new ones, please refer to the [Developer Guide](docs/DEVGUIDE.md).

## Frontend Functionality

The application includes a JavaScript file (`static/script/script.js`) that provides the following features:

- Drag and drop file upload
- File selection via input
- Date selection
- Jira comment and summary generation
- Copying generated content to clipboard
- Editing and updating system prompts
- Viewing and refreshing application logs
- Collapsible sections for better organization

For more information about the frontend implementation, please see the [Developer Guide](docs/DEVGUIDE.md).

## Developer Guide

For detailed information about the application's structure, how `main.go` works, third-party package usage, environment variables, and how to extend the project with new endpoints, please refer to the [Developer Guide](docs/DEVGUIDE.md).

The Developer Guide includes:

- An overview of `main.go` structure and functionality
- Explanation of key third-party imports
- Details on middleware and route handlers
- Instructions for adding new endpoints
- Information about the frontend implementation

## License

This project is licensed under the MIT License. See the LICENSE file for details.