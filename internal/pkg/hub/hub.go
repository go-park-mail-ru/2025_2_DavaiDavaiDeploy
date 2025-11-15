package hub

import (
	"context"
	repo "kinopoisk/internal/pkg/users/repo/pg"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Hub struct {
	connect       sync.Map
	currentOffset time.Time
	Repo          *repo.UserRepository
}

func NewHub(repo *repo.UserRepository) *Hub {
	return &Hub{
		Repo:          repo,
		currentOffset: time.Now(),
	}
}

func (h *Hub) AddClient(userID string, client *websocket.Conn) {
	h.connect.Store(client, userID)

	//постоянно чекаем, готов ли клиент что-то принимать или он отвалился
	go func() {
		for {
			_, _, err := client.NextReader()
			if err != nil {
				_ = client.Close()
				return
			}
		}
	}()

	//когда клиент ушел - удаляем его
	client.SetCloseHandler(func(code int, text string) error {
		h.connect.Delete(client)
		return nil
	})
}

// тут у нас раннер
func (h *Hub) Run(ctx context.Context) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			h.connect.Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				userID := value.(string)
				log.Print("отправили: ", userID)
				// Получаем обновление
				order, err := h.Repo.GetUpdates(ctx, uuid.FromStringOrNil(userID), h.currentOffset)
				if err != nil {
					log.Print("нет обновления: ", userID)
					return true // продолжаем работу, но не отправляем ничего
				}
				log.Print("есть обновление: ", userID)
				// Отправляем только если есть что отправить
				if err := conn.WriteJSON(order); err != nil {
					_ = conn.Close()
					log.Print(err)
					return false
				}

				return true
			})

			// Обновляем offset после проверки
			h.currentOffset = time.Now()

		case <-ctx.Done():
			return
		}
	}
}
