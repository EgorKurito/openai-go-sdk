package openai

import "net/http"

const (
	defaultOpenAIURL = "https://api.openai.com/v1"
)

type ClientConfig struct {
	apiAuthToken string

	BaseURL      string
	Organization string
	HTTPClient   *http.Client
}

func DefaultConfig(apiAuthToken string) ClientConfig {
	return ClientConfig{
		apiAuthToken: apiAuthToken,
		BaseURL:      defaultOpenAIURL,
		Organization: "",
		HTTPClient:   &http.Client{},
	}
}
