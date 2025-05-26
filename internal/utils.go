package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func LoadEnv() {
	err := godotenv.Load("/Users/nszavadin/tmpSecrets/.env")
	if err != nil {
		log.Fatal(err)
	}
}

func NewOpenAIClient() openai.Client {
	LoadEnv()

	return openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
}
