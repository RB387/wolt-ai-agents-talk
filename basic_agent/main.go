package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/RB387/wolt-ai-agents-talk/internal"
	"github.com/openai/openai-go"
)

const systemPrompt = `
You run in a loop of Thought, Action, Pause, Observation.
At the end of the loop you output an Answer.
Use Thought to describe your thoughts about the question you have been asked.
Use Action to run one of the actions available to you - then return Pause.
Observation will be the result of running those actions.

Your available actions are:
ping:
e.g. ping: wolt.com
Does a ping command and return the response time in seconds

bash:
e.g. bash: go version
Returns the result of bash command execution

web_search:
e.g. web_search: capital of Portugal
Returns json with the urls of search results

scrape:
e.g. scrape: https://www.wolt.com
Returns the content of the given URL

Example session:
Question: How many islands make up Madeira?
Thought: I should do a web search for the Madeira
Action: web_search: Madeira
Pause

You will be called again with this:
Observation: Madeira is a Portuguese island chain made up of four islands: Madeira, Porto Santo, Desertas, and Selvagens, only two of which are inhabited (Madeira and Porto Santo.) 

You then output:
Answer: Four islands
`

// Action handlers
func ping(website string) string {
	if !strings.HasPrefix(website, "https://") && !strings.HasPrefix(website, "http://") {
		website = "https://" + website
	}

	start := time.Now()
	resp, err := http.Get(website)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start).Seconds()
	return fmt.Sprintf("%.2f seconds", duration)
}

func bash(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func webSearch(query string) string {
	// Create search client for finding URLs
	searchClient := internal.NewSearchClient()

	results, err := searchClient.Search(query)
	if err != nil {
		return fmt.Sprintf("Error performing web search: %v", err)
	}

	if len(results) == 0 {
		return "No results found"
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		return fmt.Sprintf("Error marshaling results: %v", err)
	}

	return string(jsonResults)
}

func scrape(url string) string {
	scraperClient := internal.NewScraperClient()
	content, err := scraperClient.Scrape(url)
	if err != nil {
		return fmt.Sprintf("Error scraping content: %v", err)
	}

	return content
}

// queryModel sends a query to the OpenAI API and returns the response
func queryModel(client openai.Client, messages []openai.ChatCompletionMessageParamUnion) string {
	chatCompletion, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    "gpt-4o", // Use an appropriate model
		},
	)

	if err != nil {
		fmt.Printf("Error creating chat completion: %v\n", err)
		return ""
	}

	return chatCompletion.Choices[0].Message.Content
}

// runAgentLoop executes the agent's thought-action-observation loop
func runAgentLoop(client openai.Client, userQuery string, maxIter int) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(userQuery),
	}

	actionRegex := regexp.MustCompile(`^Action: (\w+): (.*)`)

	for i := 0; i < maxIter; i++ {
		fmt.Printf("Loop: %d\n", i+1)

		response := queryModel(client, messages)
		fmt.Println(response)

		// Find action in response
		lines := strings.Split(response, "\n")
		var actionFound bool

		for _, line := range lines {
			matches := actionRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				action := matches[1]
				actionInput := matches[2]

				fmt.Printf("Running %s %s\n", action, actionInput)

				var observation string
				switch action {
				case "ping":
					observation = ping(actionInput)
				case "bash":
					observation = bash(actionInput)
				case "web_search":
					observation = webSearch(actionInput)
				case "scrape":
					observation = scrape(actionInput)
				default:
					observation = "Unknown action"
				}

				fmt.Printf("Observation: %s\n", observation)

				// Add observation to messages
				messages = append(messages, openai.UserMessage(fmt.Sprintf("Observation: %s", observation)))
				actionFound = true
				break
			}
		}

		if !actionFound {
			fmt.Println("No more actions, agent is done.")
			break
		}
	}
}

func main() {
	client := internal.NewOpenAIClient()

	fmt.Println("=== Query 1: Response time for wolt.com ===")
	runAgentLoop(client, "What's the response time for wolt.com?", 5)

	fmt.Println("\n=== Query 2: Go version ===")
	runAgentLoop(client, "What version of Golang is installed on this machine?", 5)

	fmt.Println("\n=== Query 3: What is the weather in Helsinki today ===")
	runAgentLoop(client, "What is the weather in Helsinki today (in Celsius)? Also print time when the weather was checked. Peferably from accuweather", 5)
}
