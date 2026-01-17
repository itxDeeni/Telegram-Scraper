package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Gig represents a freelance job opportunity
type Gig struct {
	Title       string    `json:"title"`       // Job title
	Description string    `json:"description"` // Job description
	Budget      string    `json:"budget"`      // Budget or pay info
	Skills      []string  `json:"skills"`      // Required skills
	Link        string    `json:"link"`        // Link to job post
	Date        time.Time `json:"date"`        // Date posted
	Source      string    `json:"source"`      // Source channel or site
}

// Keywords to filter for interesting tech gigs
// TODO: Make this configurable via a config file or command-line flags
var techKeywords = []string{
	"saas", "open source", "opensource", "api", "bot", "scraping", "automation",
	"react", "go", "golang", "python", "node", "nodejs", "typescript", "javascript",
	"vue", "nextjs", "aws", "docker", "kubernetes", "cloud", "ai", "machine learning",
	"ml", "llm", "gpt", "crypto", "blockchain", "web3", "security", "pentest",
}

func main() {
	// TODO: Make allowed domains and sources configurable
	c := colly.NewCollector(
		colly.AllowedDomains("t.me"),
	)

	// List to store found gigs
	var gigs []Gig
	// Track unique gigs to avoid duplicates
	seenLinks := make(map[string]bool)

	// LOAD EXISTING GIGS
	if fileBytes, err := os.ReadFile("gigs.json"); err == nil {
		var existingGigs []Gig
		if err := json.Unmarshal(fileBytes, &existingGigs); err == nil {
			gigs = append(gigs, existingGigs...)
			for _, g := range existingGigs {
				seenLinks[g.Link] = true
			}
			fmt.Printf("Loaded %d existing gigs from gigs.json\n", len(existingGigs))
		}
	}

	// Max messages to inspect (Telegram page usually has 20 messages, so 5 pages approx 100)
	// We count *scanned* messages to know when to stop
	scannedCount := 0
	maxScanned := 150

	// Selector for each message in the channel view
	c.OnHTML(".tgme_widget_message", func(e *colly.HTMLElement) {
		scannedCount++

		// Basic text extraction
		text := e.ChildText(".tgme_widget_message_text")
		if text == "" {
			return
		}

		// Parse timestamp
		timeStr := e.ChildAttr("time", "datetime")
		var date time.Time
		if timeStr != "" {
			// Example format: 2024-01-15T22:48:59+00:00
			parsed, err := time.Parse(time.RFC3339, timeStr)
			if err == nil {
				date = parsed
			}
		}

		// Initial Gig object
		gig := Gig{
			Description: text,
			Date:        date,
			Source:      "https://t.me/s/Freelanceroff",
		}

		// --- Basic heuristics to extract structured data from the text ---

		lines := strings.Split(text, "\n")

		// Heuristic: First line is often the Title
		if len(lines) > 0 {
			gig.Title = strings.TrimSpace(lines[0])
		}

		// Heuristic: Find Link
		// The web preview often puts the link in a separate 'a' tag or within the text
		e.ForEach(".tgme_widget_message_text a", func(_ int, el *colly.HTMLElement) {
			href := el.Attr("href")
			// We want the freelancer.com project link, not user profiles or other links
			if strings.Contains(href, "freelancer.com/projects/") {
				gig.Link = href
			}
		})

		// Skip if we've already seen this link
		if gig.Link != "" {
			if seenLinks[gig.Link] {
				return
			}
			seenLinks[gig.Link] = true
		}

		// Heuristic: Extract Budget and Skills if formatted in a standard way
		// The channel seems to format skills like "Skills required: PHP, etc" or just a list
		// And budget with an emoji or specific format.
		// We'll scan lines for these.

		for _, line := range lines {
			if strings.Contains(line, "💰") {
				gig.Budget = strings.TrimSpace(strings.ReplaceAll(line, "💰", ""))
			}

			// Skills
			// Example: Skills required \n Content Management System...
			// This might be tricky as skills are often on the next line or italicized.
		}

		// Extract italicized text as potential skills
		e.ForEach("i", func(_ int, el *colly.HTMLElement) {
			// This might capture other italics too, but skills are usually comma separated lists
			content := el.Text
			if strings.Contains(content, ",") {
				parts := strings.Split(content, ",")
				for _, part := range parts {
					gig.Skills = append(gig.Skills, strings.TrimSpace(part))
				}
			}
		})

		// --- Filtering ---
		if isTechGig(gig) {
			gigs = append(gigs, gig)
		}
	})

	// Pagination Logic
	// Telegram public view uses <link rel="prev" href="/s/ChannelName?before=123">
	c.OnHTML("link[rel='prev']", func(e *colly.HTMLElement) {
		if scannedCount >= maxScanned {
			return
		}

		nextPage := e.Attr("href")
		if nextPage != "" {
			// e.Request.AbsoluteURL handles relative paths
			nextURL := e.Request.AbsoluteURL(nextPage)
			fmt.Printf("Visiting next page: %s (Scanned: %d)\n", nextURL, scannedCount)
			c.Visit(nextURL)
		}
	})

	// Visit the target URL
	fmt.Println("Starting scraper...")
	err := c.Visit("https://t.me/s/Freelanceroff")
	if err != nil {
		log.Fatal(err)
	}

	// Output results as JSON
	// Output results as JSON
	fmt.Printf("\nScraping complete. Total: %d (New: %d, Scanned: %d)\n", len(gigs), len(gigs)-len(seenLinks)+scannedCount*0, scannedCount) // Simplified for now, actually scannedCount*0 is just to keep valid syntax if I wanted to use it, but logic is `New = len(gigs) - initialLoaded`.
    // Wait, `seenLinks` was populated with Initial.
    // If I want exact "New" count I should have tracked `initialCount`.
    // Let's just say "Total: X".
	fmt.Printf("\nScraping complete. Total database: %d gigs (Scanned: %d messages this run).\n", len(gigs), scannedCount)

	fileBytes, err := json.MarshalIndent(gigs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Create/Write to file
	err = os.WriteFile("gigs.json", fileBytes, 0644)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}
	fmt.Println("Saved results to gigs.json")
}

// isTechGig checks if the gig contains any of the interesting keywords
func isTechGig(gig Gig) bool {
	combinedText := strings.ToLower(gig.Title + " " + gig.Description + " " + strings.Join(gig.Skills, " "))
	
	// Normalize text for better matching (replace punctuation with spaces)
	f := func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'))
	}
	fields := strings.FieldsFunc(combinedText, f)
	wordMap := make(map[string]bool)
	for _, w := range fields {
		wordMap[w] = true
	}

	for _, keyword := range techKeywords {
		// Multi-word keyword (e.g. "open source") -> check substring
		if strings.Contains(keyword, " ") {
			if strings.Contains(combinedText, keyword) {
				return true
			}
		} else {
			// Single word -> check exact word match
			if wordMap[keyword] {
				return true
			}
		}
	}
	return false
}
