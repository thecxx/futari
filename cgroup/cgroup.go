package cgroup

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrInvalidCommandArgs = errors.New("invalid command args")
)

type Command interface {
	Do() (err error)
}

type Creator struct {
	New      func(args []string, sayFn func(string)) Command
	template string
	explain  string
}

var (
	cmds  = make(map[string]Creator)
	mutex sync.RWMutex
)

// GetPrompt
func GetPrompt() (prompt string) {
	if len(cmds) <= 0 {
		return
	}

	prompts := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		prompts = append(prompts, fmt.Sprintf("%s\n%s", cmd.template, cmd.explain))
	}

	return strings.Join(prompts, "\n\n")
}

// RegisterCommand
func RegisterCommand(name string, newFn func(args []string, sayFn func(string)) Command, template, explain string) {
	mutex.Lock()
	defer mutex.Unlock()
	cmds[name] = Creator{New: newFn, template: template, explain: explain}
}

// QueryCommand
func QueryCommand(name string) (newFn func(args []string, sayFn func(string)) Command, supported bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	cmd, supported := cmds[name]
	if supported {
		newFn = cmd.New
	}
	return
}
