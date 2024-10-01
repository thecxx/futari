package futari

import (
	"context"
	"sync"

	"github.com/thecxx/futari/define"
	"github.com/thecxx/futari/define/types"
)

type Engine interface {
	Chat(ctx context.Context, messages []types.Message) (answer types.Message, err error)
}

// Model
type Model struct {
	engine    Engine
	prompt    []types.Message
	histories map[string][]types.Message
	mutex     sync.RWMutex
}

// NewModel
func NewModel(engine Engine, prompt types.Message) (mod *Model) {
	mod = &Model{
		engine:    engine,
		histories: make(map[string][]types.Message),
	}
	// Initial prompt
	mod.prompt = append(mod.prompt, prompt)
	return
}

// Tell
func (mod *Model) Tell(ctx context.Context, message types.Message) (answer types.Message, err error) {
	mod.mutex.Lock()
	defer mod.mutex.Unlock()

	var (
		messages []types.Message
	)

	// Can be any user
	role := message.Role
	// Must be "user"
	message.Role = define.RoleUser

	messages = append(messages, mod.prompt...)
	history, ok := mod.histories[role]
	if !ok {
		history = make([]types.Message, 0)
		mod.histories[role] = history
	}
	messages = append(messages, history...)
	messages = append(messages, message)

	// Send messages
	answer, err = mod.engine.Chat(ctx, messages)
	if err != nil {
		return
	}

	history = append(history, answer)

	mod.histories[role] = history
	return
}

// GetHistory
func (mod *Model) GetHistory(role string) (history []types.Message) {
	mod.mutex.RLock()
	defer mod.mutex.RUnlock()
	return mod.histories[role]
}

// RemoveHistory
func (mod *Model) RemoveHistory(role string, messageID uint64) {
	mod.mutex.Lock()
	defer mod.mutex.Unlock()
	// Update history
	history, ok := mod.histories[role]
	if !ok {
		return
	}
	newHistory := make([]types.Message, 0)
	for _, message := range history {
		if message.ID != messageID {
			newHistory = append(newHistory, message)
		}
	}
	mod.histories[role] = newHistory
}
