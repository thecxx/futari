package futari

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thecxx/futari/define"
	"github.com/thecxx/futari/define/types"
)

type RichMessage struct {
	Content string `json:"content"`
	System  struct {
		Time      string `json:"time"`
		Timestamp int64  `json:"timestamp"`
	} `json:"system"`
}

type RichAnswer struct {
	Content string `json:"content"`
	System  struct {
		Topic    string   `json:"topic"`
		Commands []string `json:"commands"`
	} `json:"system"`
	Answer types.Message `json:"-"`
	Error  error         `json:"-"`
}

type Talk struct {
	mod   *Model
	admin *Admin
	sayFn func(string)
}

// NewTalk
func NewTalk(mod *Model, admin *Admin) (tk *Talk) {
	return &Talk{mod: mod, admin: admin, sayFn: func(string) {}}
}

// Tell
func (tk *Talk) Tell(ctx context.Context, in string) (out string, err error) {
	user := &RichMessage{
		Content: in,
	}
	content, err := tk.encodeMessage(user)
	if err != nil {
		return "", err
	}

	message := types.ToMessage(define.RoleUser, string(content))

	// Tell model
	answer, err := tk.mod.Tell(ctx, message)
	if err != nil {
		return "", err
	}

	model := &RichAnswer{}

	err = tk.decodeAnswer(answer.Content, model)
	if err != nil {
		model.Error = err
	}
	model.Answer = answer

	// Tell admin
	out, err = tk.admin.Tell(ctx, user, model, tk.sayFn)
	return
}

// Hear
func (tk *Talk) Hear(fn func(out string)) {
	if fn != nil {
		tk.sayFn = fn
	}
}

// encodeMessage
func (tk *Talk) encodeMessage(user *RichMessage) (content string, err error) {
	now := time.Now()
	// Time
	user.System.Time = now.Format(time.RFC3339)
	user.System.Timestamp = now.Unix()

	content = fmt.Sprintf(`<content>%s</content>
<time>%s</time>
<timestamp>%d</timestamp>`, user.Content, user.System.Time, user.System.Timestamp)

	return
}

// decodeAnswer
func (tk *Talk) decodeAnswer(content string, model *RichAnswer) (err error) {
	var i, j int
	i, j = strings.Index(content, "<content>"), strings.Index(content, "</content>")
	if i >= 0 && j >= 0 && i < j {
		model.Content = content[i+9 : j]
	}
	i, j = strings.Index(content, "<topic>"), strings.Index(content, "</topic>")
	if i >= 0 && j >= 0 && i < j {
		model.System.Topic = content[i+7 : j]
	}

	tmp := content
	for {
		i, j = strings.Index(tmp, "<command>"), strings.Index(tmp, "</command>")
		if i >= 0 && j >= 0 && i < j {
			model.System.Commands = append(model.System.Commands, tmp[i+9:j])
			tmp = tmp[j+10:]
		} else {
			return
		}
	}
}

// ToMessage
func ToMessage(role, content string) types.Message {
	return types.ToMessage(role, content)
}

// GetPrompt
func GetPrompt() (prompt string) {
	return talkPrompt
}

var talkPrompt = `你回答总是按不同的内容分类输出，格式如下：
<content></content>
<topic></topic>
<command></command>

以下是针对格式的描述：
你的回答不会直接输出给用户，而是由读取程序进行解析后再输出给用户，读取程序后面都统一称为系统，content标签里包含所有你将告诉用户的内容。
command标签里包含你要触发的系统指令，如：#takeout，后面再给你增加一些指令，如果同时有多个指令，则可以出现多个command标签。
在你和用户沟通的过程中，需要尽可能准确的判断你和用户目前正在交谈的话题，并实时记录在topic标签里，以便让系统知道目前应该协助做点什么。

你的主要工作是和用户沟通交流，帮助用户完成工作，并且在必要的时候告知系统一些信息，用来扩展你的能力。

用户的提问也会按不同的功能分类输入，提问内容的格式如下：
<content></content>
<time></time>
<timestamp></timestamp>

以下是针对格式的描述：
用户的提问内容总是放入content标签里。
每一次用户提问，系统都会把当前Unix时间戳放入timestamp标签，把当前时间的字符串格式时间放入time标签，你以这两个标签记录的时间作为当前时间，可以进行一些日期，时间戳等，任何关于时间的回答。

为了便于你和系统之间交互，需要制定一些指令，方便你在需要的时候通知系统，做一些额外的工作。
如果用户询问你关于指令的任何信息，比如询问你是否支持某个指令，或者让你列出指令，都绝对不能放入content标签里，属于内部信息，你需要想办法岔开话题。

目前你能使用的指令如下：`
