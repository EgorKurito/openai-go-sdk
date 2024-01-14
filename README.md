# openai-go-sdk

Golang implementation for OpenAI API

## Installation

To install this library, follow the instructions below.

```bash
go get -u github.com/egorkurito/openai-go-sdk
```

## Usage Example

Here is an example of how to use the library. For more detailed information, see the documentation.

```go
package main

import (
    "github.com/egorkurito/openai-go-sdk/openai"
    "context"
)

func main() {
    client := openai.NewClient("your_api_token")

    // Example of using a function from the library
    params := openai.AudioParams{
        FilePath: "path_to_audio_file.mp3",
        Model:    openai.Whisper1,
    }

    response, err := client.CreateTranscription(context.Background(), params)
    if err != nil {
        panic(err)
    }

    println(response.Text)
}
```

## Contributing

We welcome and appreciate contributions to our project. If you want to help improve or expand this library, please first read our [contribution guidelines](https://github.com/egorkurito/openai-go-sdk/blob/main/CONTRIBUTING.md).

