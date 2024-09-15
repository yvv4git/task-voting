package infrastructure

import (
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
)

func TestNewSubscription(t *testing.T) {
	logger := NewDefaultLogger()
	s := NewSubscription(logger)
	if len(s.clients) != 0 {
		t.Errorf("Expected empty clients map, got %v", s.clients)
	}
}

func TestAddClient(t *testing.T) {
	logger := NewDefaultLogger()
	s := NewSubscription(logger)
	client := &websocket.Conn{}
	s.AddClient(client)
	if _, ok := s.clients[client]; !ok {
		t.Errorf("Expected client to be added to clients map")
	}
}

func TestRemoveClient(t *testing.T) {
	logger := NewDefaultLogger()
	s := NewSubscription(logger)
	client := &websocket.Conn{}
	s.AddClient(client)
	s.RemoveClient(client)
	if _, ok := s.clients[client]; ok {
		t.Errorf("Expected client to be removed from clients map")
	}
}

type mockConn struct {
	Conn         ClientConn
	writeMessage func(mt int, data []byte) error
}

func (c *mockConn) WriteMessage(mt int, data []byte) error {
	if c.writeMessage != nil {
		return c.writeMessage(mt, data)
	}
	return c.Conn.WriteMessage(mt, data)
}

func (c *mockConn) Close() error {
	return nil
}

func TestBroadcast(t *testing.T) {
	logger := NewDefaultLogger()
	subscription := NewSubscription(logger)
	client1 := &mockConn{Conn: &websocket.Conn{}}
	client2 := &mockConn{Conn: &websocket.Conn{}}
	subscription.AddClient(client1)
	subscription.AddClient(client2)

	client1.writeMessage = func(mt int, data []byte) error {
		return fmt.Errorf("mock error")
	}

	client2.writeMessage = func(mt int, data []byte) error {
		return nil
	}

	message := []byte("Hello")
	subscription.Broadcast(message)

	_, hasClient1 := subscription.clients[client1]
	if hasClient1 {
		t.Errorf("Expected client1 to be removed from clients map")
	}
}

func TestBroadcastNoClients(t *testing.T) {
	logger := NewDefaultLogger()
	s := NewSubscription(logger)
	message := []byte("Hello")
	s.Broadcast(message)
	// No error expected
}
