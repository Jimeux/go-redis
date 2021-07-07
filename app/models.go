package app

import (
	"strconv"
	"time"
)

var (
	TableThreads   = "threads"
	TableMessages  = "messages"
	TableReactions = "reactions"
	TableUsers     = "users"

	PKThreads   = "id"
	PKMessages  = "thread_id"
	PKUsers     = "id"
	PKReactions = "thread_id"

	SKMessages  = "sent_at"
	SKReactions = "sort_key"
)

type Thread struct {
	ID        string    `dynamodbav:"id"`
	Title     string    `dynamodbav:"title"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}

func NewThread(id string, title string) *Thread {
	return &Thread{ID: id, Title: title, CreatedAt: time.Now().UTC()}
}

type Message struct {
	ThreadID string `dynamodbav:"thread_id"`
	Content  string `dynamodbav:"content"`
	UserID   int64  `dynamodbav:"user_id"`
	SentAt   int64  `dynamodbav:"sent_at"`
}

func MessageID(threadID string, sentAt int64) string {
	return threadID + ":" + strconv.FormatInt(sentAt, 10)
}

func MessageIDs(msgs []*Message) []string {
	ids := make([]string, len(msgs))
	for i, msg := range msgs {
		ids[i] = MessageID(msg.ThreadID, msg.SentAt)
	}
	return ids
}

func NewMessage(threadID string, content string, userID int64) *Message {
	return &Message{
		ThreadID: threadID,
		Content:  content,
		UserID:   userID,
		SentAt:   time.Now().UTC().UnixNano() / 1000,
	}
}

type Reaction struct {
	ThreadID      string `dynamodbav:"thread_id"`
	SortKey       string `dynamodbav:"sort_key"`
	MessageSentAt int64  `dynamodbav:"message_sent_at"`
	UserID        int64  `dynamodbav:"user_id"`
	Kind          string `dynamodbav:"kind"`
}

func NewReaction(threadID string, sentAt int64, userID int64, kind string) *Reaction {
	return &Reaction{
		ThreadID:      threadID,
		SortKey:       strconv.FormatInt(sentAt, 10) + ":" + strconv.FormatInt(userID, 10) + ":" + kind,
		UserID:        userID,
		Kind:          kind,
		MessageSentAt: sentAt,
	}
}

type User struct {
	ID       string `dynamodbav:"id"`
	Username string `dynamodbav:"username"`
	Avatar   string `dynamodbav:"avatar"`
}
