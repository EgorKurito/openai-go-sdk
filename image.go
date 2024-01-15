package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	GenerateImageModelDallE2 = "dall-e-2"
	GenerateImageModelDallE3 = "dall-e-3"

	GenerateImageQualityHD       = "hd"
	GenerateImageQualityStandard = "standard"

	GenerateImageResponseFormatURL     = "url"
	GenerateImageResponseFormatB64JSON = "b64_json"

	GenerateImageSize256x256   = "256x256"
	GenerateImageSize512x512   = "512x512"
	GenerateImageSize1024x1024 = "1024x1024"
	GenerateImageSize1792x1024 = "1792x1024"
	GenerateImageSize1024x1792 = "1024x1792"

	GenerateImageStyleVivid   = "vivid"
	GenerateImageStyleNatural = "natural"
)

// GenerateImageParams - represents the request structure for API.
type GenerateImageParams struct {
	// A text description of the desired image(s).
	// The maximum length is 1000 characters for `dall-e-2` and 4000 characters for `dall-e-3`.
	Prompt string `json:"prompt,omitempty"`

	// The model to use for image generation.
	Model string `json:"model,omitempty"`

	// The number of images to generate.
	// Must be between 1 and 10. For `dall-e-3`, only `n=1` is supported.
	N int `json:"n,omitempty"`

	// The quality of the image that will be generated.
	// `hd` creates images with finer details and greater consistency across the image.
	// This param is only supported for `dall-e-3`.
	Quality string `json:"quality,omitempty"`

	// The format in which the generated images are returned.
	// Must be one of `url` or `b64_json`.
	ResponseFormat string `json:"response_format,omitempty"`

	// The size of the generated images.
	// Must be one of `256x256`, `512x512`, or `1024x1024` for `dall-e-2`. Must be one
	// of `1024x1024`, `1792x1024`, or `1024x1792` for `dall-e-3` models.
	Size string `json:"size,omitempty"`

	// The style of the generated images.
	// Must be one of `vivid` or `natural`. Vivid causes the model to lean towards
	// generating hyper-real and dramatic images. Natural causes the model to produce
	// more natural, less hyper-real looking images. This param is only supported for `dall-e-3`.
	Style string `json:"style,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices/end-user-ids).
	User string `json:"user,omitempty"`
}

// GenerateImageResponse - represents a response structure for API.
type GenerateImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL           string `json:"url,omitempty"`
		B64JSON       string `json:"b64_json,omitempty"`
		RevisedPrompt string `json:"revised_prompt,omitempty"`
	} `json:"data"`
}

// GenerateImage -  API call to generate image.
func (c *Client) GenerateImage(ctx context.Context, params GenerateImageParams) (*GenerateImageResponse, error) {
	response, err := c.callGenerateImageAPI(ctx, params, "/images/generations")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// callGenerateImageAPI - API call to an image generating endpoint.
func (c *Client) callGenerateImageAPI(
	ctx context.Context,
	params GenerateImageParams,
	endpointSuffix string,
) (*GenerateImageResponse, error) {
	var response GenerateImageResponse

	var reqBytes []byte

	reqBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getFullURL(endpointSuffix), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if err = c.sendRequest(req, &response); err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	return &response, nil
}
