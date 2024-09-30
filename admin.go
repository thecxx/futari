package futari

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/thecxx/futari/command"
	"github.com/thecxx/futari/common"
)

var (
	regexpCmd = regexp.MustCompile(`#([^\s]+)\s+(.+)`)
)

// Admin
type Admin struct {
}

// NewAdmin
func NewAdmin() (admin *Admin) {
	return new(Admin)
}

func (admin *Admin) Tell(ctx context.Context, message common.Message, reaction Reaction) (out string, err error) {
	out = reaction.Message

	// Nothing to do
	if len(reaction.System.Commands) <= 0 {
		return
	}

	for _, cmd := range reaction.System.Commands {
		submatch := regexpCmd.FindStringSubmatch(cmd)
		if len(submatch) > 0 {
			name := submatch[1]
			args := make([]string, 0)
			if len(submatch) > 2 {
				args = strings.Split(submatch[2], " ")
			}
			// Query command creator
			creator, ok := command.QueryCommand(name)
			if !ok {
				return "", fmt.Errorf("command: %s not supported", name)
			}
			// Create command
			executor := creator(args)
			// Execute command
			if err = executor.Do(); err != nil {
				return "", err
			}
		}
	}

	return
}
