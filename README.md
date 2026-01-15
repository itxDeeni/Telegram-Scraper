# Freelance Gig Scraper

A Go-based tool to scrape and filter freelance tech gigs from public Telegram channel.

## Features
- Scrapes job listings from [Freelancer | Freelance](https://t.me/s/Freelanceroff).
- Filters for tech-related keywords (SaaS, Go, React, Python, API, etc.).
- handled pagination to fetch historical data.
- Outputs results to `gigs.json`.

## Usage

1. **Install dependencies**
   ```bash
   go mod tidy
   ```

2. **Run the scraper**
   ```bash
   go run main.go
   ```

3. **View Results**
   The tool will save relevant gigs to `gigs.json` and print a summary to the console.

## Configuration
You can modify the `techKeywords` slice in `main.go` to customize the filtering logic.

## Disclaimer
This tool is for educational purposes only. It scrapes public web preview of the Freelancer Telegram channel.
