package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

// Implement marshaller interface and use Scan

type SelfMarshaller struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
	SentAt   int64  `json:"sent_at"`
}

func (m *SelfMarshaller) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
func (m *SelfMarshaller) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, m)
}

func BenchmarkSelfMarshaller(b *testing.B) {
	c := client()
	ctx := context.Background()
	d := time.Second * 10
	m := &SelfMarshaller{
		ThreadID: "thread",
		Content:  "content",
		SentAt:   12343452345,
	}
	for i := 0; i < b.N; i++ {
		key := "slf-marshal:" + strconv.Itoa(i)
		if _, err := c.Set(ctx, key, m, d).Result(); err != nil {
			b.Fatal(err)
		}
		var m2 SelfMarshaller
		if err := c.Get(ctx, key).Scan(&m2); err != nil {
			b.Fatal(err)
		}
		if m2.ThreadID == "" {
			b.Fatal("failed to parse")
		}
	}
}

// Use json.Marshal/Unmarshal manually

type NonMarshaller struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
	SentAt   int64  `json:"sent_at"`
}

func BenchmarkNonMarshaller(b *testing.B) {
	c := client()
	ctx := context.Background()
	d := time.Second * 10
	m := &NonMarshaller{
		ThreadID: "thread",
		Content:  "content",
		SentAt:   12343452345,
	}
	for i := 0; i < b.N; i++ {
		raw, err := json.Marshal(m)
		if err != nil {
			b.Fatal(err)
		}

		key := "non-marshal:" + strconv.Itoa(i)
		if _, err := c.Set(ctx, key, raw, d).Result(); err != nil {
			b.Fatal(err)
		}
		 res, err := c.Get(ctx, key).Bytes()
		 if err != nil {
			b.Fatal(err)
		}
		var m2 SelfMarshaller
		if err := json.Unmarshal(res, &m2); err != nil {
			b.Fatal(err)
		}
		if m2.ThreadID == "" {
			b.Fatal("failed to parse")
		}
	}
}

// Use protobuf

func (m *ProtoMarshaller) MarshalBinary() ([]byte, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (m *ProtoMarshaller) UnmarshalBinary(b []byte) error {
	if err := proto.Unmarshal(b, m); err != nil {
		return err
	}
	return nil
}

func BenchmarkProtoMarshaller(b *testing.B) {
	c := client()
	ctx := context.Background()
	d := time.Second * 10
	m := &ProtoMarshaller{
		ThreadId: "thread",
		Content:  "content",
		SentAt:   12343452345,
	}
	for i := 0; i < b.N; i++ {
		key := "non-marshal:" + strconv.Itoa(i)
		if _, err := c.Set(ctx, key, m, d).Result(); err != nil {
			b.Fatal(err)
		}
		var m2 ProtoMarshaller
		if err := c.Get(ctx, key).Scan(&m2); err != nil {
			b.Fatal(err)
		}
		if m2.ThreadId == "" {
			b.Fatal("failed to parse")
		}
	}
}

func client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:26379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
