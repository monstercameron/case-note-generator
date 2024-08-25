#!/bin/bash

# Check if Go 1.22 or newer is installed
if command -v go >/dev/null 2>&1; then
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "1.22" "$go_version" | sort -V | head -n1)" = "1.22" ]; then
        echo -e "\e[32mGo $go_version is installed.\e[0m"
    else
        echo -e "\e[31mGo $go_version is installed, but version 1.22 or newer is required.\e[0m"
        echo -e "\e[31mPlease upgrade Go from https://golang.org/dl/\e[0m"
        exit 1
    fi
else
    echo -e "\e[31mGo is not installed. Please install Go 1.22 or newer from https://golang.org/dl/\e[0m"
    exit 1
fi

# Setup the project
echo -e "\e[36mSetting up the project...\e[0m"

# Go up one folder from current working directory
cd ..

# Check if .env file exists and warn if values are empty
if [ -f ".env" ]; then
    echo -e "\e[32m.env file found.\e[0m"
    empty_vars=()
    while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ "$line" =~ ^[[:alnum:]_]+=$ ]]; then
            empty_vars+=("${line%=}")
        fi
    done < .env
    if [ ${#empty_vars[@]} -gt 0 ]; then
        echo -e "\e[33mWarning: The following environment variables are empty:\e[0m"
        for var in "${empty_vars[@]}"; do
            echo -e "\e[33m  - $var\e[0m"
        done
    fi
else
    echo -e "\e[33mWarning: .env file not found in the parent directory.\e[0m"
fi

# Check if static and logs folders exist
if [ ! -d "static" ]; then
    echo -e "\e[33mWarning: 'static' folder does not exist.\e[0m"
fi
if [ ! -d "logs" ]; then
    echo -e "\e[33mWarning: 'logs' folder does not exist.\e[0m"
fi

# Initialize go module (if not already initialized)
if [ ! -f "go.mod" ]; then
    go mod init myproject
    echo -e "\e[32mGo module initialized.\e[0m"
else
    echo -e "\e[33mGo module already initialized.\e[0m"
fi

# Download dependencies
echo -e "\e[36mDownloading dependencies...\e[0m"
if go mod download; then
    echo -e "\e[32mDependencies downloaded successfully.\e[0m"
else
    echo -e "\e[31mError downloading dependencies. Please check your go.mod file and internet connection.\e[0m"
fi

echo -e "\e[32mProject setup complete!\e[0m"

# Print useful commands
echo -e "\n\e[33mUseful commands:\e[0m"
echo "- Run the application: go run main.go"
echo "- Build the application: go build -o myapp main.go"
echo "- Run tests: go test ./..."
echo "- Format code: go fmt ./..."
echo "- Lint code: go vet ./..."
echo -e "\nRemember to update the .env file with your actual API key and create necessary folders if they don't exist."