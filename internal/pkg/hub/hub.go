package hub

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"kinopoisk/internal/pkg/films/repo"
	"kinopoisk/internal/pkg/utils/log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connect       sync.Map
	currentOffset time.Time
	Repo          *repo.FilmRepository
}

func (h *Hub) AddClient(client *websocket.Conn) {
	h.connect.Store(client, true)
	go func() {
		for {
			_, _, err := client.NextReader()
			if err != nil {
				h.connect.Delete(client)
				_ = client.Close()
				return
			}
		}
	}()
	client.SetCloseHandler(func(code int, text string) error {
		h.connect.Delete(client)
		return nil
	})
}

func (h *Hub) Run(ctx context.Context) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			news, ok := h.Repo.GetUpdates(ctx, h.currentOffset)
			if !ok {
				h.currentOffset = time.Now()
				continue
			}
			logger.Info("sending news to all clients", "count:", len(news))
			h.connect.Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				if conn == nil {
					h.connect.Delete(key)
					return true
				}
				if err := conn.WriteJSON(news); err != nil {
					h.connect.Delete(key)
					_ = conn.Close()
					logger.Error("failed to write news: ", err)
					return true
				}

				return true
			})
			h.currentOffset = time.Now()

		case <-ctx.Done():
			h.connect.Range(func(key, value interface{}) bool {
				if conn, ok := key.(*websocket.Conn); ok {
					_ = conn.Close()
				}
				return true
			})
			return
		}
	}
}
