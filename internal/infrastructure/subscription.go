package infrastructure

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientConn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type Subscription struct {
	logger  *slog.Logger
	clients map[ClientConn]bool
	mu      sync.Mutex
}

func NewSubscription(logger *slog.Logger) *Subscription {
	return &Subscription{
		logger:  logger,
		clients: make(map[ClientConn]bool),
	}
}

func (s *Subscription) AddClient(client ClientConn) {
	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()
}

func (s *Subscription) RemoveClient(client ClientConn) {
	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
}

func (s *Subscription) Broadcast(message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			client.Close()
			delete(s.clients, client)
			s.logger.Error("Client disconnected", slog.String("error", err.Error()))
		}
	}
}
