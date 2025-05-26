package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RB387/wolt-ai-agents-talk/internal"
	swarmgo "github.com/prathyushnallamothu/swarmgo"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

// Tool to search the web
func searchWeb(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	query, ok := args["query"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "Error: query parameter is required and must be a string",
			Success: false,
		}
	}

	query = strings.ReplaceAll(query, `"`, "")
	fmt.Println("Searching the web for:", query)

	searchClient := internal.NewSearchClient()
	results, err := searchClient.Search(query)
	if err != nil {
		return swarmgo.Result{
			Data: fmt.Sprintf("Error searching the web: %v", err),
			Success: false,
		}
	}

	if len(results) == 0 {
		return swarmgo.Result{
			Data: "No results found",
			Success: false,
		}
	}

	var resultText strings.Builder
	for i, result := range results {
		if i > 1 { // Limit to 2 results
			break
		}
		resultText.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, result.URL, result.URL))
	}

	return swarmgo.Result{
		Data: resultText.String(),
		Success: true,
	}
}

// Tool to scrape content from a URL
func scrapeUrl(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	url, ok := args["url"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "Error: url parameter is required and must be a string",
			Success: false,
		}
	}

	fmt.Println("Scraping URL:", url)

	url = strings.ReplaceAll(url, `"`, "")
	scraperClient := internal.NewScraperClient()
	content, err := scraperClient.Scrape(url)
	if err != nil {
		return swarmgo.Result{
			Data: fmt.Sprintf("Error scraping URL: %v", err),
			Success: false,
		}
	}

	return swarmgo.Result{
		Data: content,
		Success: true,
	}
}

// Tool to manage files
func manageFiles(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	action, ok := args["action"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "Error: action parameter is required and must be a string (read, write, or list)",
			Success: false,
		}
	}

	path, ok := args["path"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "Error: path parameter is required and must be a string",
			Success: false,
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return swarmgo.Result{
			Data: fmt.Sprintf("Error getting current directory: %v", err),
			Success: false,
		}
	}
	path = filepath.Join(dir, path)

	fmt.Println("Managing files with action:", action, "and path:", path)

	switch action {
	case "read":
		content, err := os.ReadFile(path)
		if err != nil {
			return swarmgo.Result{
				Data: fmt.Sprintf("Error reading file: %v", err),
				Success: false,
			}
		}
		return swarmgo.Result{
			Data: string(content),
			Success: true,
		}
	case "write":
		content, ok := args["content"].(string)
		if !ok {
			return swarmgo.Result{
				Data: "Error: content parameter is required for write operation",
				Success: false,
			}
		}
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return swarmgo.Result{
				Data: fmt.Sprintf("Error writing file: %v", err),
				Success: false,
			}
		}
		return swarmgo.Result{
			Data: fmt.Sprintf("File %s written successfully", path),
			Success: true,
		}
	case "list":
		files, err := os.ReadDir(path)
		if err != nil {
			return swarmgo.Result{
				Data: fmt.Sprintf("Error listing directory: %v", err),
				Success: false,
			}
		}
		var result strings.Builder
		for _, file := range files {
			result.WriteString(file.Name() + "\n")
		}
		return swarmgo.Result{
			Data: result.String(),
			Success: true,
		}
	default:
		return swarmgo.Result{
			Data: fmt.Sprintf("Unknown action: %s. Must be read, write, or list", action),
			Success: false,
		}
	}
}

// Tool to get human input
func getHumanInput(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	question, ok := args["question"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "Error: question parameter is required",
			Success: false,
		}
	}

	fmt.Printf("\n🧠 Human input needed: %s\n", question)
	fmt.Print("Your response: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	return swarmgo.Result{
		Data: response,
		Success: true,
	}
}

func main() {
	internal.LoadEnv()

	workflow := swarmgo.NewWorkflow(os.Getenv("OPENAI_API_KEY"), llm.OpenAI, swarmgo.SupervisorWorkflow)
	workflow.SetCycleHandling(swarmgo.ContinueOnCycle)

	// Define supervisor functions
	supervisorFunctions := []swarmgo.AgentFunction{
		{
			Name:        "getHumanInput",
			Description: "Ask a human for input",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"question": map[string]interface{}{
						"type":        "string",
						"description": "The question to ask the human",
					},
				},
				"required": []string{"question"},
			},
			Function: getHumanInput,
		},
	}

	// Define writer functions
	writerFunctions := []swarmgo.AgentFunction{
		{
			Name:        "manageFiles",
			Description: "Read, write, and list files",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "The action to perform: read, write, list or mkdir",
						"enum":        []string{"read", "write", "list", "mkdir"},
					},
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The file or directory path",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content to write (only for write action)",
					},
				},
				"required": []string{"action", "path"},
			},
			Function: manageFiles,
		},
	}

	// Define scraper functions
	scraperFunctions := []swarmgo.AgentFunction{
		{
			Name:        "searchWeb",
			Description: "Search the web and return the json with the urls of search results",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
				},
				"required": []string{"query"},
			},
			Function: searchWeb,
		},
		{
			Name:        "scrapeUrl",
			Description: "Fetch the content of a URL",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to scrape",
					},
				},
				"required": []string{"url"},
			},
			Function: scrapeUrl,
		},
	}

	// Create supervisor agent
	supervisorAgent := &swarmgo.Agent{
		Name: "supervisor",
		Instructions: `
		You are a supervisor tasked with managing a conversation between the following workers: 
		[scraper, writer]. 
		Given the following user request, respond with the worker to act next. 
		Each worker will perform a task and respond with their results and status. 
		You can ask human for input anytime.

		Scrapper is responsible for finding and extracting information from the web.
		Writer is responsible for creating a comprehensive report.

		Scrapper can search the web and return the json with the urls of search results.
		Using the urls, scrapper can scrape the information from the web with separate call.

		Writer aggregates the information from the scrapper and writes a comprehensive report.

		You should not do anything else.
		When finished, respond with FINISH.`,
		Functions: supervisorFunctions,
		Model:     "gpt-4.1",
	}

	// Create writer agent
	writerAgent := &swarmgo.Agent{
		Name: "writer",
		Instructions: `You are the writer agent responsible for creating a comprehensive report.
Your role is to:
1. Draft the report content based on information provided by the supervisor
2. Ensure consistent tone and style throughout the document
3. Organize content with proper structure, headings, and formatting
4. Write Final Report to to file in the Markdown format`,
		Functions: writerFunctions,
		Model:     "gpt-4.1",
	}

	// Create scraper agent
	scraperAgent := &swarmgo.Agent{
		Name: "scraper",
		Instructions: `You are the scraper agent responsible for finding and extracting information from the web.
Your role is to:
1. SEARCH for information using the searchWeb function to find relevant URLs
1.1 Simply return list of urls

2. SCRAPE specific URLs from those search results using the scrapeUrl function

IMPORTANT: After scraping, extract content of the page. Return clean text, not raw HTML.`,
		Functions: scraperFunctions,
		Model:     "gpt-4.1",
	}

	// Add agents to teams
	workflow.AddAgentToTeam(supervisorAgent, swarmgo.SupervisorTeam)
	workflow.AddAgentToTeam(writerAgent, swarmgo.DocumentTeam)
	workflow.AddAgentToTeam(scraperAgent, swarmgo.ResearchTeam)

	// Set supervisor as team leader
	if err := workflow.SetTeamLeader(supervisorAgent.Name, swarmgo.SupervisorTeam); err != nil {
		log.Fatal("Error setting team leader:", err)
	}

	// Connect agents
	workflow.ConnectAgents(supervisorAgent.Name, writerAgent.Name)
	workflow.ConnectAgents(supervisorAgent.Name, scraperAgent.Name)
	workflow.ConnectAgents(writerAgent.Name, supervisorAgent.Name)
	workflow.ConnectAgents(scraperAgent.Name, supervisorAgent.Name)

	// Define the user request
	userPrompt := "I need a comprehensive report on something that is happening in the world. You can ask human for any clarifications."

	printWorkflowStart()

	result, err := workflow.Execute(supervisorAgent.Name, userPrompt)
	if err != nil {
		log.Fatal("Error executing workflow:", err)
	}

	printWorkflowSummary(*result)

	for _, step := range result.Steps {
		printStepResult(step)
	}
}


//
// PRINT UTILS
//
func printOutputs(output []llm.Message) {
	for _, msg := range output {
		switch msg.Role {
		case llm.RoleUser:
			fmt.Printf("\033[92m[User]\033[0m: %s\n", msg.Content)
		case llm.RoleAssistant:
			name := msg.Name
			if name == "" {
				name = "Assistant"
			}
			fmt.Printf("\033[94m[%s]\033[0m: %s\n", name, msg.Content)
		case llm.RoleFunction, "tool":
			fmt.Printf("\033[95m[Function Result]\033[0m: %s\n", msg.Content)
		}
	}
}


func printStepResult(step swarmgo.StepResult) {
	fmt.Println("\n\033[96mDetailed Step Results\033[0m")

	fmt.Printf("\n\033[95mStep %d Results:\033[0m\n", step.StepNumber)
	fmt.Printf("Agent: %s\n", step.AgentName)
	fmt.Printf("Duration: %v\n", step.EndTime.Sub(step.StartTime))
	if step.Error != nil {
		fmt.Printf("\033[91mError: %v\033[0m\n", step.Error)
		return
	}

	fmt.Println("Output:")
	printOutputs(step.Output)

	if step.NextAgent != "" {
		fmt.Printf("\nNext Agent: %s\n", step.NextAgent)
	}
	fmt.Println("-----------------------------------------")
}

func printWorkflowSummary(result swarmgo.WorkflowResult) {
	fmt.Printf("\n\033[96mWorkflow Summary\033[0m\n")
	fmt.Printf("Total Duration: %v\n", result.EndTime.Sub(result.StartTime))
	fmt.Printf("Total Steps: %d\n", len(result.Steps))
}

func printWorkflowStart() {
	fmt.Println("\n\033[96mStarting Report Generation Workflow\033[0m")
	fmt.Println("================================")
}
