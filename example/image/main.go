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

	resp, err := client.GenerateImage(
		context.Background(),
		openai.GenerateImageParams{
			Prompt:         "Pixel programming cat with a hat",
			Model:          openai.GenerateImageModelDallE2,
			Size:           openai.GenerateImageSize512x512,
			ResponseFormat: openai.GenerateImageResponseFormatURL,
			N:              1,
		},
	)
	if err != nil {
		log.Fatalf("Transcription error: %v\n", err)

		return
	}

	log.Printf(resp.Data[0].URL)
}
