package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Jimeux/go-redis/cache"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type CreateThreadReq struct {
	Title string `json:"title"`
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	var req CreateThreadReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, err)
		return
	}
	msg, err := h.svc.CreateThread(r.Context(), req.Title)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, msg)
}

type FindMessagesRes struct {
	Messages cache.Messages `json:"messages"`
	CacheHit bool           `json:"cache_hit"`
}

func (h *Handler) FindMessages(w http.ResponseWriter, r *http.Request) {
	threadID := r.URL.Query().Get("thread_id")
	min, _ := strconv.ParseInt(r.URL.Query().Get("min"), 10, 64)
	max, _ := strconv.ParseInt(r.URL.Query().Get("max"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	msgs, hit, err := h.svc.FindMessages(r.Context(), threadID, min, max, limit)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, FindMessagesRes{
		Messages: msgs,
		CacheHit: hit,
	})
}

type SendMessageReq struct {
	ThreadID string `json:"thread_id"`
	UserID   int64  `json:"user_id"`
	Content  string `json:"content"`
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, err)
		return
	}
	msg, err := h.svc.SendMessage(r.Context(), req.ThreadID, req.Content, req.UserID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, msg)
}

type SendReactionReq struct {
	ThreadID string `json:"thread_id"`
	UserID   int64  `json:"user_id"`
	Kind     string `json:"kind"`
	SentAt   int64  `json:"sent_at"`
}

func (h *Handler) SendReactions(w http.ResponseWriter, r *http.Request) {
	var req SendReactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, err)
		return
	}
	res, err := h.svc.SendReaction(r.Context(), req.ThreadID, req.Kind, req.UserID, req.SentAt)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	h.writeJSON(w, r, res)
}

// helper methods

func (h *Handler) writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.writeError(w, r, err)
	}
}

func (*Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(r)
}
