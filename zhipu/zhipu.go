package zhipu

import (
	"context"
	"net/http"

	"github.com/thecxx/futari/define"
	"github.com/thecxx/futari/define/types"
	"github.com/yankeguo/zhipu"
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

func NewZhipu(model string, opts ...zhipu.ClientOption) (z *Zhipu, err error) {
	client, err := zhipu.NewClient(opts...)
	if err != nil {
		return
	}
	z = &Zhipu{model: model, client: client}
	return
}

// SendMessages implements futari.Model.
func (z *Zhipu) SendMessages(ctx context.Context, messages []types.Message) (answer types.Message, err error) {
	service := z.client.ChatCompletion(z.model)

	for _, v := range messages {
		service = service.AddMessage(zhipu.ChatCompletionMessage{Role: v.Role, Content: v.Content})
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return answer, err
	}

	return types.ToMessage(define.RoleAssistant, resp.Choices[0].Message.Content), nil
}
