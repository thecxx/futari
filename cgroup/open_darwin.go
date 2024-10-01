package cgroup

import (
	"fmt"
	"os/exec"
)

var openTemplate = "[#open 应用名]"
var openExplain = `打开指令会通知系统你要触发的指令是打开应用，指令包含固定的#open和一个应用名参数，格式固定，不要随便修改。
如果你判断用户可能需要打开应用，你可以询问用户是否需要帮助打开应用，如果用户确认需要，注意只有当用户明确同意你帮助打开时，你才能使用该指令，如果不确认就不要使用该指令。`

func init() {
	RegisterCommand("open", NewOpen, openTemplate, openExplain)
}

type Open struct {
	name        string
	application string
	args        []string
	say         func(string)
}

// NewOpen
func NewOpen(args []string, sayFn func(string)) Command {
	c := &Open{name: "启动器", args: args, say: sayFn}
	if len(args) >= 1 {
		c.application = args[0]
	}
	return c
}

// Do
func (c *Open) Do() (err error) {
	if c.application == "" {
		return ErrInvalidCommandArgs
	}
	cmd := exec.Command("open", fmt.Sprintf("/Applications/%s.app", c.application))
	return cmd.Start()
}
