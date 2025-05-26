package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/RB387/wolt-ai-agents-talk/internal"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

// Custom tool implementations
type PingTool struct{}

func (p PingTool) Name() string {
	return "ping"
}

func (p PingTool) Description() string {
	return "Ping a website and get the response time"
}

func (p PingTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("Pinging:", input)
	if !strings.HasPrefix(input, "https://") && !strings.HasPrefix(input, "http://") {
		input = "https://" + input
	}

	start := time.Now()
	resp, err := http.Get(input)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	duration := time.Since(start).Seconds()
	return fmt.Sprintf("%.2f seconds", duration), nil
}

type BashTool struct{}

func (b BashTool) Name() string {
	return "bash"
}

func (b BashTool) Description() string {
	return "Execute a bash command and return the output"
}

func (b BashTool) Call(ctx context.Context, command string) (string, error) {
	fmt.Println("Executing bash command:", command)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

type WebSearchTool struct{}

func (w WebSearchTool) Name() string {
	return "web_search"
}

func (w WebSearchTool) Description() string {
	return "Search the web and return the json with the urls of search results"
}

func (w WebSearchTool) Call(ctx context.Context, query string) (string, error) {
	query = strings.ReplaceAll(query, `"`, "")
	fmt.Println("Searching the web for:", query)
	// Create search client for finding URLs
	searchClient := internal.NewSearchClient()

	results, err := searchClient.Search(query)
	if err != nil {
		return "", nil
	}

	if len(results) == 0 {
		return "", nil
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		return "", nil
	}

	return string(jsonResults), nil
}

type ScrapeTool struct{}

func (s ScrapeTool) Name() string {
	return "scrape"
}

func (s ScrapeTool) Description() string {
	return "Fetch the content of a URL"
}

func (s ScrapeTool) Call(ctx context.Context, url string) (string, error) {
	url = strings.ReplaceAll(url, `"`, "")
	fmt.Println("Scraping:", url)
	scraperClient := internal.NewScraperClient()
	content, err := scraperClient.Scrape(url)
	if err != nil {
		return "", err
	}

	return content, nil
}

func run() error {
	internal.LoadEnv()

	llm, err := openai.New(openai.WithToken(os.Getenv("OPENAI_API_KEY")), openai.WithModel("gpt-4.1"))
	if err != nil {
		return fmt.Errorf("error initializing OpenAI client: %w", err)
	}

	// Set up tools
	agentTools := []tools.Tool{
		PingTool{},
		BashTool{},
		WebSearchTool{},
		ScrapeTool{},
	}

	// Create agent and executor
	agent := agents.NewOneShotAgent(llm, agentTools, agents.WithMaxIterations(5))
	executor := agents.NewExecutor(agent)

	// Run the first query
	fmt.Println("=== Query 1: Response time for wolt.com ===")
	result1, err := chains.Run(context.Background(), executor, "What's the response time for wolt.com?")
	if err != nil {
		return fmt.Errorf("error executing query 1: %w", err)
	}
	fmt.Println("Result:", result1)

	// Run the second query
	fmt.Println("\n=== Query 2: Go version ===")
	result2, err := chains.Run(context.Background(), executor, "What version of Golang is installed on this machine?")
	if err != nil {
		return fmt.Errorf("error executing query 2: %w", err)
	}
	fmt.Println("Result:", result2)

	// Run the third query
	fmt.Println("\n=== Query 3: Weather in Helsinki ===")
	result3, err := chains.Run(context.Background(), executor, "What is the weather in Helsinki today (in Celsius)? Peferably from accuweather ")
	if err != nil {
		return fmt.Errorf("error executing query 3: %w", err)
	}
	fmt.Println("Result:", result3)

	return nil
}


func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}