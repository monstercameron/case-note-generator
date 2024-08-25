<#
.SYNOPSIS
    Build and run script for Go project.
.DESCRIPTION
    This script checks for Go installation, verifies project structure,
    downloads dependencies, builds, and runs the current Go project.
.NOTES
    File Name      : setup.ps1
    Author         : Your Name
    Prerequisite   : PowerShell 5.0 or later
    Copyright 2023 - Your Company
#>

# Output formatting functions
function Write-CheckboxStatus($message, $status) {
    $checkbox = if ($status) { "[V]" } else { "[ ]" }
    $color = if ($status) { "Green" } else { "Yellow" }
    Write-Host "$checkbox $message" -ForegroundColor $color
}

function Write-Warning($message) {
    Write-Host "[WARN] $message" -ForegroundColor Yellow
}

function Write-Error($message) {
    Write-Host "[ERROR] $message" -ForegroundColor Red
}

function Write-Info($message) {
    Write-Host "[INFO] $message" -ForegroundColor Cyan
}

# Check if Go 1.22 or newer is installed
function Test-GoVersion {
    Write-Info "Checking Go version..."
    try {
        $goVersion = (go version)
        if ($goVersion -match "go version go(\d+\.\d+)") {
            $version = [version]$Matches[1]
            $minVersion = [version]"1.22"
            if ($version -ge $minVersion) {
                Write-CheckboxStatus "Go $($version.ToString()) is installed." $true
                return $true
            } else {
                Write-CheckboxStatus "Go $($version.ToString()) is installed. Version 1.22 or newer is required." $false
                return $false
            }
        } else {
            Write-CheckboxStatus "Unable to parse Go version." $false
            return $false
        }
    } catch {
        Write-CheckboxStatus "Go is not installed or not in PATH." $false
        return $false
    }
}

# Check .env file and warn about empty values
function Test-EnvFile {
    Write-Info "Checking .env file..."
    $envPath = Join-Path (Split-Path $PSScriptRoot -Parent) ".env"
    if (Test-Path $envPath) {
        Write-CheckboxStatus ".env file exists." $true
        $emptyValues = Get-Content $envPath | Where-Object { $_ -match "^(\w+)=\s*$" } | ForEach-Object { $Matches[1] }
        if ($emptyValues) {
            Write-Warning "The following environment variables are empty:"
            $emptyValues | ForEach-Object { Write-Warning "  - $_" }
        } else {
            Write-Info "All environment variables are set."
        }
    } else {
        Write-CheckboxStatus ".env file exists." $false
    }
}

# Check if required folders exist
function Test-RequiredFolders {
    Write-Info "Checking required folders..."
    $projectRoot = Split-Path $PSScriptRoot -Parent
    @("static", "logs") | ForEach-Object {
        $folderPath = Join-Path $projectRoot $_
        $folderExists = Test-Path $folderPath
        Write-CheckboxStatus "'$_' folder exists." $folderExists
    }
}

# Download dependencies
function Get-GoDependencies {
    Write-Info "Downloading dependencies..."
    $projectRoot = Split-Path $PSScriptRoot -Parent
    Push-Location $projectRoot
    $output = go mod download 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-CheckboxStatus "Dependencies downloaded successfully." $true
    } else {
        Write-CheckboxStatus "Error downloading dependencies." $false
        Write-Error "Error details: $output"
    }
    Pop-Location
}

# Build the project
function Build-Project {
    Write-Info "Building the project..."
    $projectRoot = Split-Path $PSScriptRoot -Parent
    Push-Location $projectRoot
    $output = go build -o myapp 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-CheckboxStatus "Project built successfully." $true
        return $true
    } else {
        Write-CheckboxStatus "Error building the project." $false
        Write-Error "Error details: $output"
        return $false
    }
    Pop-Location
}

# Run the project
function Run-Project {
    Write-Info "Running the project..."
    $projectRoot = Split-Path $PSScriptRoot -Parent
    $appPath = Join-Path $projectRoot "myapp"
    if (Test-Path $appPath) {
        Start-Process -FilePath $appPath -NoNewWindow
        Write-CheckboxStatus "Project started." $true
    } else {
        Write-CheckboxStatus "Project executable not found." $false
        Write-Error "The 'myapp' executable was not found. Make sure the build process completed successfully."
    }
}

# Main execution
function Main {
    Write-Host "`n=== Go Project Build and Run Script ===`n" -ForegroundColor Magenta

    if (-not (Test-GoVersion)) {
        Write-Error "Go installation check failed. Please install Go 1.22 or newer from https://golang.org/dl/"
        exit 1
    }

    Test-EnvFile
    Test-RequiredFolders
    Get-GoDependencies
    $buildSuccess = Build-Project
    if ($buildSuccess) {
        Run-Project
    }

    Write-Host "`n=== Project Setup Complete ===`n" -ForegroundColor Magenta

    if ($buildSuccess) {
        Write-Host "The project has been built and is now running." -ForegroundColor Yellow
        Write-Host "You can stop the application by closing its window or using Ctrl+C in its console." -ForegroundColor Yellow
    } else {
        Write-Host "The project build failed. Please review the errors above and try again." -ForegroundColor Yellow
    }
    Write-Host ""
    Write-Host "Remember to review any warnings or errors reported during the process." -ForegroundColor Yellow
}

# Run the main function
Main