package futari

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/thecxx/futari/cgroup"
)

var (
	regexpCmd = regexp.MustCompile(`^#([^\s]+)\s+(.+)$`)
)

// Admin
type Admin struct {
}

// NewAdmin
func NewAdmin() (admin *Admin) {
	return new(Admin)
}

func (admin *Admin) Tell(ctx context.Context, user *RichMessage, model *RichAnswer, sayFn func(out string)) (out string, err error) {
	if model.Error != nil {
		err = model.Error
		return
	}

	out = model.Content

	// Nothing to do
	if len(model.System.Commands) <= 0 {
		return
	}

	for _, cmd := range model.System.Commands {
		submatch := regexpCmd.FindStringSubmatch(cmd)
		if len(submatch) > 0 {
			name := submatch[1]
			args := make([]string, 0)
			if len(submatch) > 2 {
				args = strings.Split(submatch[2], " ")
			}

			// Query command
			newCommand, ok := cgroup.QueryCommand(name)
			if !ok {
				return "", fmt.Errorf("command: %s not supported", name)
			}

			// Create command
			c := newCommand(args, sayFn)
			// Execute command
			if err = c.Do(); err != nil {
				return "", err
			}
		}
	}

	return
}
