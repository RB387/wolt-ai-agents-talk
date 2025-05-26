package main

import (
	"context"
	"fmt"

	"github.com/RB387/wolt-ai-agents-talk/internal"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)


func queryModel(client openai.Client, prompt string) string {
	chatCompletion, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
				openai.SystemMessage(""),
			},
			Model: shared.ChatModelGPT4o,
		},
	)

	if err != nil {
		fmt.Printf("Error creating chat completion: %v\n", err)
		return ""
	}

	return chatCompletion.Choices[0].Message.Content
}


func main() {
	client := internal.NewOpenAIClient()

	// First query
	result1 := queryModel(client, "What's the response time for wolt.com?")
	fmt.Println(result1)
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")

	// Second query
	result2 := queryModel(client, "What version of Golang is installed on this machine?")
	fmt.Println(result2)
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")

	result3 := queryModel(client, "What's the weather in Helsinki today?")
	fmt.Println(result3)
}
