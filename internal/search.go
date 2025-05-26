package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// SearchResult represents a single search result
type SearchResult struct {
	URL string `json:"url"`
}

// SearchResponse represents the response from the search API
type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// SearchClient provides search functionality
type SearchClient struct {
	apiHost string
	baseURL string
	client  *http.Client
}

// NewSearchClient creates a new search client
func NewSearchClient() *SearchClient {
	LoadEnv()

	return &SearchClient{
		apiHost: "duckduckgo8.p.rapidapi.com",
		baseURL: "https://duckduckgo8.p.rapidapi.com/",
		client:  &http.Client{},
	}
}

// Search performs a search query and returns results
func (s *SearchClient) Search(query string) ([]SearchResult, error) {
	// URL encode the query
	encodedQuery := url.QueryEscape(query)
	requestURL := fmt.Sprintf("%s?q=%s", s.baseURL, encodedQuery)

	// Create request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", s.apiHost)

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var searchResponse SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return searchResponse.Results, nil
}
