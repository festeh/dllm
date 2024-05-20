package dllm

import (
	"net/http"
	"net/url"
)

type OpenaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIData struct {
	Model       string          `json:"model"`
	Messages    []OpenaiMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	Stream      bool            `json:"stream"`
}

type OpenAI struct {
	authToken string
	client    *http.Client
}

type OpenAIQuery struct {
	Messages []OpenaiMessage `json:"messages"`
}

func NewOpenAI(authToken string) *OpenAI {
	return &OpenAI{authToken, &http.Client{}}
}

func (o *OpenAI) Name() string {
	return "OpenAI"
}

func (o *OpenAI) CompletionURL() *url.URL {
	return ParseUrlYolo("https://api.openai.com/v1/chat/completions")
}


func (o *OpenAI) GetStream(body []byte, writer http.ResponseWriter) (*Stream, error) {
	return NewStream(body, writer, o)
}

func (a *OpenAI) LoadQuery(body []byte) (*OpenAIQuery, error) {
	query := OpenAIQuery{}
	err := LoadQuery(body, &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

func (o *OpenAI) CreateData(query *OpenAIQuery) *OpenAIData {
	return &OpenAIData{
		Model:       "gpt-4",
		Messages:    query.Messages,
		Temperature: 0.1,
		MaxTokens:   500,
		Stream:      true,
	}
}

func (o *OpenAI) addHeaders(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+o.authToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
}

func (o *OpenAI) do(request *http.Request) (*http.Response, error) {
	return o.client.Do(request)
}

func (o *OpenAI) GetWriterCallback() func([]byte) {
	return func(body []byte) {
		// fmt.Fprintf(w, "data: %s\n\n", body)
	}
}

