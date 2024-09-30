package common

import (
	"sync/atomic"
)

type Message struct {
	ID      uint64
	Role    string
	Content string
}

var (
	messageID uint64
)

// ToMessage
func ToMessage(role, content string) Message {
	return Message{ID: atomic.AddUint64(&messageID, 1), Role: role, Content: content}
}
