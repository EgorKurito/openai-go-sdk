package main

import (
	"context"
	"log"
	"os"

	"github.com/egorkurito/openai-go-sdk"
	"github.com/joho/godotenv"
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
		log.Fatalf("Transcription error: %v\n", err)

		return
	}

	log.Printf(resp.Text)
}
