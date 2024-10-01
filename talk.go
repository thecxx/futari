package futari

import (
	"context"
	"encoding/json"

	"github.com/thecxx/futari/define/types"
)

type Reaction struct {
	Message string        `json:"message"`
	System  System        `json:"system"`
	Answer  types.Message `json:"-"`
	Error   error         `json:"-"`
}

type System struct {
	Topic    string   `json:"topic"`
	Commands []string `json:"commands"`
}

type Talk struct {
	mod   *Model
	admin *Admin
}

// NewTalk
func NewTalk(mod *Model, admin *Admin) (tk *Talk) {
	return &Talk{mod: mod, admin: admin}
}

// Tell
func (tk *Talk) Tell(ctx context.Context, message types.Message) (out string, err error) {
	// Tell model
	answer, err := tk.mod.Tell(ctx, message)
	if err != nil {
		return "", err
	}

	var reaction Reaction

	err = tk.parseAnswer(answer.Content, &reaction)
	if err != nil {
		reaction.Error = err
	}
	reaction.Answer = answer

	// Tell admin
	out, err = tk.admin.Tell(ctx, message, reaction)
	return
}

// parseAnswer
func (tk *Talk) parseAnswer(content string, reaction *Reaction) (err error) {
	return json.Unmarshal([]byte(content), reaction)
}

// ToMessage
func ToMessage(role, content string) types.Message {
	return types.ToMessage(role, content)
}
