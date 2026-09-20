// Package realtime = broadcast skor via Redis pub/sub + WebSocket.
// Redis jadi jembatan antar-instance backend (siap multi-instance, "gantian server").
package realtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// ScoreMessage = event skor yang disiarkan ke channel match.
type ScoreMessage struct {
	MatchID string          `json:"match_id"`
	Score   json.RawMessage `json:"score"`
}

const channelPrefix = "kampiun:match:"

// Hub menampung koneksi WS per match & meneruskan pesan dari Redis.
type Hub struct {
	rdb  *redis.Client
	log  *slog.Logger
	mu   sync.RWMutex
	subs map[string]*redis.PubSub   // matchID -> subscription
	conns map[string]map[*wsConn]struct{} // matchID -> koneksi
	upgrader websocket.Upgrader
}

type wsConn struct{ send chan []byte }

// NewHub membuat hub; terhubung ke Redis bila redisAddr diisi.
func NewHub(redisAddr string, log *slog.Logger) *Hub {
	h := &Hub{
		log: log,
		subs: map[string]*redis.PubSub{},
		conns: map[string]map[*wsConn]struct{}{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}
	if redisAddr != "" {
		h.rdb = redis.NewClient(&redis.Options{Addr: redisAddr})
	}
	return h
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, s := range h.subs {
		_ = s.Close()
	}
	if h.rdb != nil {
		_ = h.rdb.Close()
	}
}

// channel utk match.
func channelFor(matchID string) string { return channelPrefix + matchID }

// Publish mengirim skor ke channel Redis (dipanggil score service setelah tap).
func (h *Hub) Publish(ctx context.Context, matchID string, scoreJson []byte) error {
	if h.rdb == nil {
		return nil // Redis nonaktif, skip
	}
	msg, _ := json.Marshal(ScoreMessage{MatchID: matchID, Score: scoreJson})
	return h.rdb.Publish(ctx, channelFor(matchID), msg).Err()
}

// BroadcastToLocal meneruskan ke koneksi WS lokal (fallback bila tanpa Redis).
func (h *Hub) BroadcastToLocal(matchID string, scoreJson []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.conns[matchID] {
		select {
		case c.send <- scoreJson:
		default: // antre penuh, drop
		}
	}
}

// ServeWS mengupgrade koneksi & menambahkannya ke match + listen Redis.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, matchID string) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &wsConn{send: make(chan []byte, 32)}
	h.mu.Lock()
	if h.conns[matchID] == nil {
		h.conns[matchID] = map[*wsConn]struct{}{}
	}
	h.conns[matchID][c] = struct{}{}
	h.mu.Unlock()

	// pastikan subscribe Redis utk match ini ada
	h.ensureSub(matchID)

	// writer
	go func() {
		defer conn.Close()
		for msg := range c.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// reader: abaikan pesan masuk, tunggu close
	defer func() {
		h.mu.Lock()
		delete(h.conns[matchID], c)
		h.mu.Unlock()
		close(c.send)
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

// ensureSub memastikan satu goroutine subscribe Redis utk match tsb jalan.
func (h *Hub) ensureSub(matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.subs[matchID]; ok {
		return
	}
	if h.rdb == nil {
		// tanpa redis: tak ada sub; broadcast lokal via Publish fallback
		h.subs[matchID] = nil
		return
	}
	sub := h.rdb.Subscribe(context.Background(), channelFor(matchID))
	h.subs[matchID] = sub
	go func() {
		ch := sub.Channel()
		for msg := range ch {
			var m ScoreMessage
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				continue
			}
			h.BroadcastToLocal(m.MatchID, m.Score)
		}
	}()
}
