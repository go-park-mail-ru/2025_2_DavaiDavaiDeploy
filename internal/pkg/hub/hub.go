package hub

import (
	"context"
	"sync"
	"time"

	"kinopoisk/internal/models"
	authrepo "kinopoisk/internal/pkg/auth/repo"
	"kinopoisk/internal/pkg/films/repo"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Hub struct {
	connect        sync.Map
	currentOffset  time.Time
	passwordOffset time.Time
	Repo           *repo.FilmRepository
	AuthRepo       *authrepo.AuthRepository
}

func (h *Hub) AddClient(userID uuid.UUID, client *websocket.Conn) {
	h.connect.Store(client, userID)
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
	t := time.NewTicker(5 * time.Second)
	t2 := time.NewTicker(30 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.connect.Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				userID := value.(uuid.UUID)
				news, ok := h.Repo.GetUpdates(ctx, userID, h.currentOffset)
				if !ok || len(news) == 0 {
					return true
				}

				go h.sendNewsToClient(conn, news)

				return true
			})

			h.currentOffset = time.Now()

		case <-t2.C:
			h.connect.Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				userID := value.(uuid.UUID)
				news, _ := h.AuthRepo.GetPasswordUpdates(ctx, userID, h.currentOffset)
				if news {
					conn.WriteJSON("Password found")
				}

				return true
			})

			h.passwordOffset = time.Now()

		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) sendNewsToClient(conn *websocket.Conn, news []models.News) {
	for _, newsItem := range news {
		if err := conn.WriteJSON(newsItem); err != nil {
			h.connect.Delete(conn)
			conn.Close()
			return
		}
	}
}
