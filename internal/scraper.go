package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ScrapeRequest represents the request body for the scraping API
type ScrapeRequest struct {
	URL string `json:"url"`
}

// ScrapeResponse represents the response from the scraping API
type ScrapeResponse struct {
	Body string `json:"body"`
}

// ScraperClient provides web scraping functionality
type ScraperClient struct {
	apiHost string
	baseURL string
	client  *http.Client
}

// NewScraperClient creates a new scraper client
func NewScraperClient() *ScraperClient {
	LoadEnv()

	return &ScraperClient{
		apiHost: "scrapeninja.p.rapidapi.com",
		baseURL: "https://scrapeninja.p.rapidapi.com/scrape",
		client:  &http.Client{},
	}
}

// Scrape fetches the content of the given URL
func (s *ScraperClient) Scrape(targetURL string) (string, error) {
	// Create request body
	reqBody := ScrapeRequest{
		URL: targetURL,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", s.apiHost)
	req.Header.Add("Content-Type", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var scrapeResp ScrapeResponse
	if err := json.Unmarshal(body, &scrapeResp); err != nil {
		// If parsing fails, just return the raw body
		return string(body), nil
	}

	return scrapeResp.Body, nil
}
