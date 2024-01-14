package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	Whisper1 = "whisper-1"
)

type AudioParams struct {
	// The audio file object (not file name) to transcribe, in one of these formats:
	// flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
	FilePath string

	// ID of the model to use. Only `whisper-1` is currently available.
	Model string

	// The language of the input audio.
	// Supplying the input language in [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) format will improve accuracy and latency.
	Language string

	// An optional text to guide the model's style or continue a previous audio segment.
	// The [prompt](https://platform.openai.com/docs/guides/speech-to-text/prompting) should match the audio language.
	Prompt string

	// The format of the transcript output, in one of these options: `json`, `text`, `srt`, `verbose_json`, or `vtt`.
	ResponseFormat string

	// The sampling temperature, between 0 and 1.
	// Higher values like 0.8 will make the output more random, while lower values like
	// 0.2 will make it more focused and deterministic. If set to 0, the model will use [log probability](https://en.wikipedia.org/wiki/Log_probability) to automatically increase the temperature until certain thresholds are hit.
	Temperature float32
}

// AudioResponse represents a response structure for audio API.
type AudioResponse struct {
	Text string `json:"text"`
}

// CreateTranscription — API call to create a transcription. Returns transcribed text.
func (c *Client) CreateTranscription(ctx context.Context, params AudioParams) (*AudioResponse, error) {
	response, err := c.callAudioAPI(ctx, params, "transcriptions")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateTranslation — API call to translate audio into English.
func (c *Client) CreateTranslation(ctx context.Context, params AudioParams) (*AudioResponse, error) {
	response, err := c.callAudioAPI(ctx, params, "translations")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// callAudioAPI — API call to an audio endpoint.
func (c *Client) callAudioAPI(ctx context.Context, params AudioParams, endpointSuffix string) (*AudioResponse, error) {
	var response AudioResponse

	var formBody bytes.Buffer
	w := multipart.NewWriter(&formBody)

	if err := audioMultipartForm(params, w); err != nil {
		return nil, err
	}

	urlSuffix := fmt.Sprintf("/audio/%s", endpointSuffix)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getFullURL(urlSuffix), &formBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	if err = c.sendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// audioMultipartForm creates a form with audio file contents and the name of the model to use for audio processing.
func audioMultipartForm(params AudioParams, w *multipart.Writer) error {
	f, err := os.Open(params.FilePath)
	if err != nil {
		return fmt.Errorf("opening audio file: %w", err)
	}
	defer f.Close()

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}

	if _, err = io.Copy(fw, f); err != nil {
		return fmt.Errorf("reading from opened audio file: %w", err)
	}

	fw, err = w.CreateFormField("model")
	if err != nil {
		return fmt.Errorf("creating form field: %w", err)
	}

	modelName := bytes.NewReader([]byte(params.Model))
	if _, err = io.Copy(fw, modelName); err != nil {
		return fmt.Errorf("writing model name: %w", err)
	}

	// Create a form field for the prompt (if provided)
	if params.Prompt != "" {
		fw, err = w.CreateFormField("prompt")
		if err != nil {
			return fmt.Errorf("creating form field: %w", err)
		}

		prompt := bytes.NewReader([]byte(params.Prompt))
		if _, err = io.Copy(fw, prompt); err != nil {
			return fmt.Errorf("writing prompt: %w", err)
		}
	}

	// Create a form field for the temperature (if provided)
	if params.Temperature != 0 {
		fw, err = w.CreateFormField("temperature")
		if err != nil {
			return fmt.Errorf("creating form field: %w", err)
		}

		temperature := bytes.NewReader([]byte(fmt.Sprintf("%.2f", params.Temperature)))
		if _, err = io.Copy(fw, temperature); err != nil {
			return fmt.Errorf("writing temperature: %w", err)
		}
	}

	// Create a form field for the language (if provided)
	if params.Language != "" {
		fw, err = w.CreateFormField("language")
		if err != nil {
			return fmt.Errorf("creating form field: %w", err)
		}

		language := bytes.NewReader([]byte(params.Language))
		if _, err = io.Copy(fw, language); err != nil {
			return fmt.Errorf("writing language: %w", err)
		}
	}

	// Close the multipart writer
	w.Close()

	return nil
}
