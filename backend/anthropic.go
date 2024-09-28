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
		query.Parameters.Model = "claude-3-5-sonnet-20240620"
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

type MessageStart struct {
	Type    string `json:"type"`
	Message struct {
		Id      string   `json:"id"`
		Type    string   `json:"type"`
		Role    string   `json:"role"`
		Content []string `json:"content"`
		Model   string   `json:"model"`
		Usage   struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

type Ping struct {
	Type string `json:"type"`
}

type ContentBlockStart struct {
	Type         string `json:"type"`
	Index        int    `json:"index"`
	ContentBlock struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block"`
}

type ContentBlockDelta struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta struct {
		Text string `json:"text"`
	} `json:"delta"`
}

type ContentBlockStop struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type MessageStop struct {
	Type string `json:"type"`
}

type MessageDelta struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason string `json:"stop_reason"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type Error struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (a *Anthropic) GetWriterCallback() func([]byte) ([]byte, bool) {
	event := []byte("event")
	messageStart := []byte("message_start")
	contentBlockStart := []byte("content_block_start")
	contentBlockDelta := []byte("content_block_delta")
	contentBlockStop := []byte("content_block_stop")
	messageStop := []byte("message_stop")
	messageDelta := []byte("message_delta")
	errorType := []byte("error")
	var dataType []byte

	data := []byte("data")
	return func(chunk []byte) ([]byte, bool) {
		if bytes.Equal(chunk[:5], event) {
			dataType = chunk[7:]
		}
		if bytes.Equal(chunk[:4], data) {
			contents := chunk[6:]
			var err error = nil
			finished := false
			output := []byte{}
			switch {
			case bytes.Equal(dataType, messageStart):
				start := &MessageStart{}
				err = json.Unmarshal(contents, start)
			case bytes.Equal(dataType, []byte("ping")):
				ping := &Ping{}
				err = json.Unmarshal(contents, ping)
			case bytes.Equal(dataType, contentBlockStart):
				start := &ContentBlockStart{}
				err = json.Unmarshal(contents, start)
			case bytes.Equal(dataType, contentBlockDelta):
				delta := &ContentBlockDelta{}
				err = json.Unmarshal(contents, delta)
				if err == nil {
					output = []byte(delta.Delta.Text)
				}
			case bytes.Equal(dataType, contentBlockStop):
				delta := &ContentBlockStop{}
				err = json.Unmarshal(contents, delta)
			case bytes.Equal(dataType, messageStop):
				delta := &MessageStop{}
				err = json.Unmarshal(contents, delta)
				finished = true
			case bytes.Equal(dataType, messageDelta):
				delta := &MessageDelta{}
				err = json.Unmarshal(contents, delta)
			case bytes.Equal(dataType, errorType):
				delta := &Error{}
				err = json.Unmarshal(contents, delta)
			default:
				errorMessage := fmt.Sprintf("Unknown data type: %s", contents)
				return []byte(errorMessage), false
			}
			if err != nil {
				errorMessage := fmt.Sprintf("Error unmarshalling %s", err)
				return []byte(errorMessage), false
			}
			return output, finished
		}
		return nil, false
	}

}
