package vyperclientgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestVyperWebsocketClient_Connect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	client := NewVyperWebsocketClient("test-api-key")
	client.BaseURL = u

	err := client.Connect(TokenEvents)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	if client.Conn == nil {
		t.Fatal("Connection is nil after successful connect")
	}

	if client.CurrentFeedType != TokenEvents {
		t.Fatalf("CurrentFeedType is %v, expected %v", client.CurrentFeedType, TokenEvents)
	}
}

func TestVyperWebsocketClient_Subscribe(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	client := NewVyperWebsocketClient("test-api-key")
	client.BaseURL = u

	err := client.Connect(TokenEvents)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	subscriptionMessage := TokenSubscriptionMessage{
		Action: Subscribe,
		Types:  []SubscriptionType{PumpfunTokens},
	}

	err = client.Subscribe(TokenEvents, subscriptionMessage)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
}

func TestVyperWebsocketClient_Listen(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade error: %v", err)
			return
		}
		defer c.Close()

		tokenPair := TokenPair{
			MarketId: "test-market",
			Name:     "Test Token",
			Symbol:   "TEST",
		}

		data, _ := json.Marshal(tokenPair)
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			t.Logf("Write error: %v", err)
			return
		}

		err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			t.Logf("Close error: %v", err)
			return
		}
	}))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	client := NewVyperWebsocketClient("test-api-key")
	client.BaseURL = u

	err := client.Connect(TokenEvents)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	receivedMessage := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	client.SetMessageHandler(func(data interface{}) {
		receivedMessage <- data
	})

	go func() {
		err := client.Listen()
		if err != nil && err.Error() != "websocket: close 1000 (normal)" {
			errorChan <- fmt.Errorf("Listen error: %v", err)
		} else {
			errorChan <- nil
		}
	}()

	select {
	case msg := <-receivedMessage:
		tokenPair, ok := msg.(*TokenPair)
		if !ok {
			t.Fatalf("Received message is not a TokenPair")
		}
		if tokenPair.MarketId != "test-market" || tokenPair.Name != "Test Token" || tokenPair.Symbol != "TEST" {
			t.Fatalf("Received unexpected TokenPair: %+v", tokenPair)
		}
	case err := <-errorChan:
		if err != nil {
			t.Fatalf("Error from Listen goroutine: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for message")
	}
}

func TestVyperWebsocketClient_Disconnect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	client := NewVyperWebsocketClient("test-api-key")
	client.BaseURL = u

	err := client.Connect(TokenEvents)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	if client.Conn != nil {
		t.Fatal("Connection is not nil after disconnect")
	}
}
