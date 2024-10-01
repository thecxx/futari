package futari

import (
	"context"
	"sync/atomic"

	"github.com/thecxx/futari/define/types"
)

type Message struct {
	ID      uint64
	Role    string
	Content string
}

var (
	messageID uint64
)

// NewMessage
func NewMessage(role, content string) Message {
	return Message{ID: atomic.AddUint64(&messageID, 1), Role: role, Content: content}
}

type Engine interface {
	SendMessages(ctx context.Context, messages []types.Message) (answer types.Message, err error)
}

// Model
type Model struct {
	engine  Engine
	prompt  []types.Message
	history []types.Message
}

// NewModel
func NewModel(engine Engine, prompt types.Message) (mod *Model) {
	return &Model{engine: engine, prompt: []types.Message{prompt}, history: make([]types.Message, 0)}
}

// Tell
func (mod *Model) Tell(ctx context.Context, message types.Message) (answer types.Message, err error) {
	var (
		messages []types.Message
	)

	messages = append(messages, mod.prompt...)
	messages = append(messages, mod.history...)
	messages = append(messages, message)

	// Send messages
	answer, err = mod.engine.SendMessages(ctx, messages)
	if err != nil {
		return
	}

	mod.history = append(mod.history, answer)
	return
}
