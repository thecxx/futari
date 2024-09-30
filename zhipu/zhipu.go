package zhipu

import (
	"context"

	"github.com/thecxx/futari/common"
	"github.com/yankeguo/zhipu"
)

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"
	Model         = "glm-4-0520"
)

type Zhipu struct {
	model  string
	client *zhipu.Client
}

func NewZhipu(model, apiKey string) (z *Zhipu, err error) {
	client, err := zhipu.NewClient(zhipu.WithAPIKey(apiKey))
	if err != nil {
		return
	}
	z = &Zhipu{model: model, client: client}
	return
}

// SendMessages implements futari.Model.
func (z *Zhipu) SendMessages(ctx context.Context, messages []common.Message) (answer common.Message, err error) {
	service := z.client.ChatCompletion(z.model)

	for _, v := range messages {
		service = service.AddMessage(zhipu.ChatCompletionMessage{Role: v.Role, Content: v.Content})
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return answer, err
	}

	return common.ToMessage(RoleAssistant, resp.Choices[0].Message.Content), nil
}
