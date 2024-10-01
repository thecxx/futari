package cgroup

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrInvalidCommandArgs = errors.New("invalid command args")
)

type Command interface {
	Do() (err error)
}

type Creator struct {
	New      func(args []string) Command
	template string
	explain  string
}

var (
	cmds  = make(map[string]Creator)
	mutex sync.RWMutex
)

// QueryPrompts
func QueryPrompts() (prompts map[string]string) {
	if len(cmds) > 0 {
		prompts = make(map[string]string)
	}
	for name, cmd := range cmds {
		prompts[name] = fmt.Sprintf("指令: %s\n%s", cmd.template, cmd.explain)
	}
	return
}

// RegisterCommand
func RegisterCommand(name string, newFn func(args []string) Command, template, explain string) {
	mutex.Lock()
	defer mutex.Unlock()
	cmds[name] = Creator{New: newFn, template: template, explain: explain}
}

// QueryCommand
func QueryCommand(name string) (newFn func(args []string) Command, supported bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	cmd, supported := cmds[name]
	if supported {
		newFn = cmd.New
	}
	return
}
