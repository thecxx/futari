package zhipu

import (
	"context"
	"errors"
	"net/http"

	"github.com/thecxx/futari/define/types"
	"github.com/yankeguo/zhipu"
)

var (
	ErrInvalidChoices = errors.New("invalid choices")
)

const (
	ModelGLM_4_0520 = "glm-4-0520"
)

type ClientOption = zhipu.ClientOption

// WithAPIKey set the api key of the client
func WithAPIKey(apiKey string) ClientOption {
	return zhipu.WithAPIKey(apiKey)
}

// WithBaseURL set the base url of the client
func WithBaseURL(baseURL string) ClientOption {
	return zhipu.WithBaseURL(baseURL)
}

// WithHTTPClient set the http client of the client
func WithHTTPClient(client *http.Client) ClientOption {
	return zhipu.WithHTTPClient(client)
}

// WithDebug set the debug mode of the client
func WithDebug(debug bool) ClientOption {
	return zhipu.WithDebug(debug)
}

type Zhipu struct {
	model  string
	client *zhipu.Client
}

func NewZhipu(model string, opts ...zhipu.ClientOption) (zp *Zhipu, err error) {
	client, err := zhipu.NewClient(opts...)
	if err != nil {
		return
	}
	zp = &Zhipu{model: model, client: client}
	return
}

// Chat implements futari.Engine.
func (zp *Zhipu) Chat(ctx context.Context, messages []types.Message) (answer types.Message, err error) {
	service := zp.client.ChatCompletion(zp.model)

	for _, v := range messages {
		message := zhipu.ChatCompletionMessage{
			Role:    v.Role,
			Content: v.Content,
		}
		service = service.AddMessage(message)
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return answer, err
	}

	if len(resp.Choices) <= 0 {
		return answer, ErrInvalidChoices
	}

	choice := resp.Choices[0]

	return types.ToMessage(choice.Message.Role, choice.Message.Content), nil
}
