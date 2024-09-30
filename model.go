package futari

import (
	"context"
	"sync/atomic"

	"github.com/thecxx/futari/common"
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
	SendMessages(ctx context.Context, messages []common.Message) (answer common.Message, err error)
}

// Model
type Model struct {
	engine  Engine
	prompt  []common.Message
	history []common.Message
}

// NewModel
func NewModel(engine Engine, prompt common.Message) (mod *Model) {
	return &Model{engine: engine, prompt: []common.Message{prompt}, history: make([]common.Message, 0)}
}

// Tell
func (mod *Model) Tell(ctx context.Context, message common.Message) (answer common.Message, err error) {
	var (
		messages []common.Message
	)

	// fmt.Printf("talk: %v\n", message)
	// defer func() {
	// 	fmt.Printf("answer: %v\n", answer)
	// }()

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
