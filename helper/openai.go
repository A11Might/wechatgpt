package helper

import (
	"context"
	"fmt"
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

func (oai *OpenAI) Chat(ctx context.Context, content string) (string, error) {
	oai.messages = append(oai.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})

	if len(oai.messages) == oai.contextSize {
		oai.messages = oai.messages[1:]
	}

	fmt.Printf("got:%+v\n", oai.messages)

	resp, err := oai.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: oai.messages,
		},
	)

	if err != nil {
		return "", err
	}
	reply := resp.Choices[0].Message.Content
	oai.messages = append(oai.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})
	return reply, nil
}

func (oai *OpenAI) STT(ctx context.Context, filePath string) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filePath,
	}

	resp, err := oai.client.CreateTranscription(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
