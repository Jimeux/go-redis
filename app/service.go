package app

import (
	"context"
	"errors"
	"log"

	"github.com/Jimeux/go-redis/cache"
)

type Service struct {
	dbRepo    *Repository
	cacheRepo *cache.Repository
}

func NewService(repo *Repository, cacheRepo *cache.Repository) *Service {
	return &Service{dbRepo: repo, cacheRepo: cacheRepo}
}

func (s *Service) CreateThread(ctx context.Context, title string) (*Thread, error) {
	t, err := s.dbRepo.SaveThread(ctx, title)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) SendMessage(ctx context.Context, threadID, content string, userID int64) (*cache.Message, error) {
	m, err := s.dbRepo.SaveMessage(ctx, threadID, content, userID)
	if err != nil {
		return nil, err
	}
	cachedMsg, err := s.cacheRepo.Save(ctx, threadID, content, userID, m.SentAt)
	if err != nil {
		return nil, err
	}
	return cachedMsg, nil
}

// SendReaction adds a new reaction to the message identified by threadID+sentAt.
// Note: It's hard to guarantee the result of FindReactions will be up-to-date
// when multiple users react simultaneously (race condition). If TTL is <=20s,
// can we accept this level of eventual consistency?
//
// Alternative approach: Completely delete the thread cache on reaction. Would the
// timing of the next read be any better than this to avoid race conditions? Probably not.
func (s *Service) SendReaction(ctx context.Context, threadID, kind string, userID, sentAt int64) (*cache.Message, error) {
	react, err := s.dbRepo.SaveReaction(ctx, threadID, sentAt, userID, kind)
	if err != nil {
		return nil, err
	}

	msgs, err := s.dbRepo.FindMessages(ctx, threadID, sentAt, sentAt, 1)
	if err != nil {
		return nil, err
	}
	if len(msgs) != 1 {
		return nil, errors.New("invalid message")
	}

	reactions, err := s.dbRepo.FindReactions(ctx, threadID, sentAt, sentAt)
	if err != nil {
		return nil, err
	}
	// reactions is unlikely to be up-to-date due to write/read latency,
	// so add react to the slice if necessary.
	reactions = addReactionIfNotContains(react, reactions)

	cachedMsg := mapMessageToCachedMessage(msgs[0], reactions)
	if err := s.cacheRepo.Update(ctx, cachedMsg); err != nil {
		return nil, err
	}
	return cachedMsg, nil
}

func (s *Service) FindMessages(ctx context.Context, threadID string, min, max, limit int64) ([]*cache.Message, bool, error) {
	cached, err := s.cacheRepo.Find(ctx, threadID, min, max, limit)
	if err != nil {
		log.Println(err)
	}

	// Return cached if non-empty AND (non-zero) max is included.
	// Lack of non-zero max could mean the cache is outdated.
	if len(cached) != 0 && (max == 0 || cached.Contains(max)) {
		return cached, true, nil
	}

	// Get from DB
	msgs, err := s.dbRepo.FindMessages(ctx, threadID, min, max, limit)
	if err != nil {
		return nil, false, err
	}
	reactions, err := s.dbRepo.FindReactions(ctx, threadID, min, max)
	if err != nil {
		return nil, false, err
	}

	// convert to response
	res := mapReactionsToMessages(reactions, msgs)

	// save to cache
	if err := s.cacheRepo.SaveAll(ctx, threadID, res); err != nil {
		log.Println(err)
	}

	return res, false, nil
}

// helpers

func mapReactionsToMessages(rr []*Reaction, mm []*Message) []*cache.Message {
	reactionMap := make(map[int64][]*cache.Reaction, len(mm))
	for _, r := range rr {
		reactionMap[r.MessageSentAt] = append(reactionMap[r.MessageSentAt], &cache.Reaction{
			UserID: r.UserID,
			Kind:   r.Kind,
		})
	}

	msgs := make([]*cache.Message, len(mm))
	for i, m := range mm {
		cm := &cache.Message{
			ThreadID:  m.ThreadID,
			Content:   m.Content,
			SentAt:    m.SentAt,
			UserID:    m.UserID,
			Reactions: []*cache.Reaction{},
		}
		if reactions, ok := reactionMap[m.SentAt]; ok {
			cm.Reactions = reactions
		}
		msgs[i] = cm
	}
	return msgs
}

func mapMessageToCachedMessage(m *Message, rr []*Reaction) *cache.Message {
	reactions := make([]*cache.Reaction, len(rr))
	for i, r := range rr {
		reactions[i] = &cache.Reaction{
			UserID: r.UserID,
			Kind:   r.Kind,
		}
	}
	return &cache.Message{
		ThreadID:  m.ThreadID,
		UserID:    m.UserID,
		Content:   m.Content,
		SentAt:    m.SentAt,
		Reactions: reactions,
	}
}

func addReactionIfNotContains(r *Reaction, rr []*Reaction) []*Reaction {
	for _, react := range rr {
		if react.ThreadID == r.ThreadID && react.SortKey == r.SortKey {
			return rr
		}
	}
	return append(rr, r)
}
