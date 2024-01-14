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
	// Supplying the input language in [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
	// format will improve accuracy and latency.
	Language string

	// An optional text to guide the model's style or continue a previous audio segment.
	// The [prompt](https://platform.openai.com/docs/guides/speech-to-text/prompting) should match the audio language.
	Prompt string

	// The format of the transcript output, in one of these options: `json`, `text`, `srt`, `verbose_json`, or `vtt`.
	ResponseFormat string

	// The sampling temperature, between 0 and 1.
	// Higher values like 0.8 will make the output more random, while lower values like
	// 0.2 will make it more focused and deterministic. If set to 0, the model will use
	// [log probability](https://en.wikipedia.org/wiki/Log_probability)
	// to automatically increase the temperature until certain thresholds are hit.
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
	writer := multipart.NewWriter(&formBody)

	if err := audioMultipartForm(params, writer); err != nil {
		return nil, err
	}

	urlSuffix := fmt.Sprintf("/audio/%s", endpointSuffix)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getFullURL(urlSuffix), &formBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	if err = c.sendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// audioMultipartForm creates a form with audio file contents and the name of the model to use for audio processing.
func audioMultipartForm(params AudioParams, w *multipart.Writer) error {
	if err := addFileToForm(w, "file", params.FilePath); err != nil {
		return err
	}

	fields := map[string]string{
		"model":       params.Model,
		"prompt":      params.Prompt,
		"temperature": fmt.Sprintf("%.2f", params.Temperature),
		"language":    params.Language,
	}

	for key, value := range fields {
		if value != "" {
			if err := addFieldToForm(w, key, value); err != nil {
				return err
			}
		}
	}

	w.Close()

	return nil
}

func addFileToForm(w *multipart.Writer, fieldName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	fw, err := w.CreateFormFile(fieldName, file.Name())
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}

	_, err = io.Copy(fw, file)

	return fmt.Errorf("copying file: %w", err)
}

func addFieldToForm(w *multipart.Writer, fieldName, fieldValue string) error {
	fw, err := w.CreateFormField(fieldName)
	if err != nil {
		return fmt.Errorf("creating form field: %w", err)
	}

	_, err = fw.Write([]byte(fieldValue))

	return fmt.Errorf("writing form field: %w", err)
}
