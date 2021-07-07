package cache

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func setKey(threadID string) string {
	return "thread:" + threadID
}

// Update updates a value in the set by score.
// Note: Members are unique in a sorted set, NOT scores.
// It's possible to have multiple members with the same score.
func (r *Repository) Update(ctx context.Context, m *Message) error {
	min, max := minMaxVals(m.SentAt, m.SentAt)

	// Use a TX to delete by score (SentAt) and then re-add m with
	// that same score. (Simulates SET by score.)
	if _, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if _, err := pipe.ZRemRangeByScore(ctx, setKey(m.ThreadID), min, max).Result(); err != nil {
			return err
		}
		if _, err := pipe.ZAdd(ctx, setKey(m.ThreadID), &redis.Z{
			Score:  float64(m.SentAt),
			Member: m,
		}).Result(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Save(ctx context.Context, threadID, content string, userID, sentAt int64) (*Message, error) {
	m := NewMessage(threadID, content, userID, sentAt)
	if _, err := r.client.ZAdd(ctx, setKey(m.ThreadID), &redis.Z{
		Score:  float64(m.SentAt),
		Member: m,
	}).Result(); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) SaveAll(ctx context.Context, threadID string, mm []*Message) error {
	zz := make([]*redis.Z, len(mm))
	for i, m := range mm {
		zz[i] = &redis.Z{
			Score:  float64(m.SentAt),
			Member: m,
		}
	}
	if _, err := r.client.ZAdd(ctx, setKey(threadID), zz...).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Find(ctx context.Context, threadID string, min, max, limit int64) (Messages, error) {
	maxVal, minVal := minMaxVals(max, min)
	res := r.client.ZRevRangeByScore(ctx, setKey(threadID), &redis.ZRangeBy{
		Min:    minVal,
		Max:    maxVal,
		Offset: 0,
		Count:  limit,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var msgs Messages
	if err := res.ScanSlice(&msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func minMaxVals(max int64, min int64) (string, string) {
	var maxVal = "+inf"
	if max != 0 {
		maxVal = strconv.FormatInt(max, 10)
	}
	var minVal = "-inf"
	if min != 0 {
		minVal = strconv.FormatInt(min, 10)
	}
	return maxVal, minVal
}
