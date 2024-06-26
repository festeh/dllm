package dllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicData struct {
	Model       string             `json:"model"`
	Messages    []AnthropicMessage `json:"messages"`
	Temperature float32            `json:"temperature"`
	MaxTokens   int                `json:"max_tokens"`
	Stream      bool               `json:"stream"`
	System      string             `json:"system"`
}

type Anthropic struct {
	key    string
	client *http.Client
}

func NewAnthropic() (*Anthropic, error) {
	keyName := "ANTHROPIC_API_KEY"
	key, ok := os.LookupEnv(keyName)
	if !ok {
		return nil, fmt.Errorf("Environment variable %s not set", keyName)
	}
	return &Anthropic{key, &http.Client{}}, nil
}

func (a *Anthropic) Name() string {
	return "Anthropic"
}

func (a *Anthropic) addHeaders(request *http.Request) {
	request.Header.Set("x-api-key", a.key)
	request.Header.Set("anthropic-version", "2023-06-01")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-Control", "no-cache")
}

func (a *Anthropic) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.anthropic.com/v1/messages")
}

func (a *Anthropic) CreateData(query *Query) ([]byte, error) {
	messages := make([]AnthropicMessage, len(query.Messages)-1)
	// skip the first message, which is the system message
	for i, message := range query.Messages {
		if i == 0 {
			continue
		}
		messages[i-1] = AnthropicMessage{message.Role, message.Content}
	}
	if query.Parameters.Model == "" {
		query.Parameters.Model = "claude-3-opus-20240229"
	}
	params := query.Parameters
	log.Info().Msgf("Creating data with model %s, temperature %f, max tokens %d", params.Model, *params.Temperature, *params.MaxTokens)
	data := &AnthropicData{
		Model:       params.Model,
		Messages:    messages,
		Temperature: *params.Temperature,
		MaxTokens:   *params.MaxTokens,
		Stream:      true,
		System:      query.Messages[0].Content,
	}
	return json.Marshal(data)
}

func (a *Anthropic) do(request *http.Request) (*http.Response, error) {
	return a.client.Do(request)
}

func (a *Anthropic) GetWriterCallback() func([]byte) ([]byte, bool) {
	eventHeader := []byte("event")
	stopSignal := []byte("message_stop")
	deltaSignal := []byte("content_block_delta")
	isDelta := false
	return func(chunk []byte) ([]byte, bool) {
		log.Debug().Msgf("Received chunk: %s", chunk)
		if bytes.Equal(chunk[:5], eventHeader) {
			if bytes.Index(chunk, stopSignal) != -1 {
				return []byte{}, true
			}
			if bytes.Index(chunk, deltaSignal) != -1 {
				isDelta = true
			}
			return []byte{}, false
		}
		if isDelta {
			data := chunk[6:]
			var delta struct {
				Delta struct {
					Text string `json:"text"`
				} `json:"delta"`
			}
			if err := json.Unmarshal([]byte(data), &delta); err != nil {
				log.Printf("error unmarshalling delta: %s", err)
				return nil, true
			}
			return []byte(delta.Delta.Text), false
		}
		return nil, false
	}
}
