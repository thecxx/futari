package futari

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/thecxx/futari/cgroup"
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

func (admin *Admin) Tell(ctx context.Context, req *RichMessage, resp *RichAnswer) (out string, err error) {
	if resp.Error != nil {
		err = resp.Error
		return
	}

	out = resp.Content

	// Nothing to do
	if len(resp.System.Commands) <= 0 {
		return
	}

	for _, cmd := range resp.System.Commands {
		submatch := regexpCmd.FindStringSubmatch(cmd)
		if len(submatch) > 0 {
			name := submatch[1]
			args := make([]string, 0)
			if len(submatch) > 2 {
				args = strings.Split(submatch[2], " ")
			}
			// Query command creator
			creator, ok := cgroup.QueryCommand(name)
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
