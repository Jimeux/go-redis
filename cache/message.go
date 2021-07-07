package cache

import (
	"encoding/json"
)

type Reaction struct {
	UserID int64  `json:"user_id"`
	Kind   string `json:"kind"`
}

type Messages []*Message

func (mm Messages) Contains(sentAt int64) bool {
	for _, m := range mm {
		if m.SentAt == sentAt {
			return true
		}
	}
	return false
}

type Message struct {
	ThreadID  string      `json:"thread_id"`
	UserID    int64       `json:"user_id"`
	Content   string      `json:"content"`
	SentAt    int64       `json:"sent_at"`
	Reactions []*Reaction `json:"reactions"`
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
func (m *Message) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, m)
}

func NewMessage(threadID, content string, userID, sentAt int64) *Message {
	return &Message{
		ThreadID: threadID,
		UserID:   userID,
		Content:  content,
		SentAt:   sentAt,
		Reactions: []*Reaction{},
	}
}
