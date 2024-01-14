package main

import (
	"context"
	"fmt"
	"github.com/egorkurito/openai-go-sdk"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %s", err.Error())
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateTranscription(
		context.Background(),
		openai.AudioParams{
			Model:    openai.Whisper1,
			FilePath: "test_voice.m4a",
		},
	)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}

	fmt.Println(resp.Text)
}
