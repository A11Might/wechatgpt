package helper

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client      *openai.Client
	messages    []openai.ChatCompletionMessage
	contextSize int
}

func NewOpenAI(contextSize int) *OpenAI {
	config := openai.DefaultConfig(DefaultConfig.OpenAIKey)
	proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	return &OpenAI{
		client:   openai.NewClientWithConfig(config),
		messages: make([]openai.ChatCompletionMessage, 0, contextSize),
	}
}

func (o *OpenAI) Chat(content string) (string, error) {
	o.messages = append(o.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})

	if len(o.messages) == o.contextSize {
		o.messages = o.messages[1:]
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: o.messages,
		},
	)

	if err != nil {
		return "", err
	}
	reply := resp.Choices[0].Message.Content
	o.messages = append(o.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})
	return reply, nil
}

var DefaultOpenAI *OpenAI
