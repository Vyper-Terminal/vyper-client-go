package vyperclientgo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type FeedType string

const (
	TokenEvents     FeedType = "token-events"
	MigrationEvents FeedType = "migration-events"
	WalletEvents    FeedType = "wallet-events"
)

type SubscriptionMessageType string

const (
	Subscribe   SubscriptionMessageType = "subscribe"
	Unsubscribe SubscriptionMessageType = "unsubscribe"
)

type SubscriptionType string

const (
	PumpfunTokens     SubscriptionType = "PumpfunTokens"
	RaydiumAmmTokens  SubscriptionType = "RaydiumAmmTokens"
	RaydiumCpmmTokens SubscriptionType = "RaydiumCpmmTokens"
	RaydiumClmmTokens SubscriptionType = "RaydiumClmmTokens"
)

type TokenSubscriptionMessage struct {
	Action SubscriptionMessageType `json:"action"`
	Types  []SubscriptionType      `json:"types"`
}

type WalletSubscriptionMessage struct {
	Action  SubscriptionMessageType `json:"action"`
	Wallets []string                `json:"wallets"`
}

type MessageHandler func(interface{})

type VyperWebsocketClient struct {
	BaseURL         string
	ApiKey          string
	Conn            *websocket.Conn
	MessageHandler  MessageHandler
	CurrentFeedType FeedType
	mu              sync.Mutex
}

func NewVyperWebsocketClient(apiKey string) *VyperWebsocketClient {
	return &VyperWebsocketClient{
		BaseURL: "wss://api.vyper.trade/api/v1/ws",
		ApiKey:  apiKey,
	}
}

func (c *VyperWebsocketClient) Connect(feedType FeedType) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn != nil {
		return fmt.Errorf("already connected")
	}

	u, err := url.Parse(fmt.Sprintf("%s/%s", c.BaseURL, feedType))
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("apiKey", c.ApiKey)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &VyperWebsocketError{
			Message:        fmt.Sprintf("Failed to connect: %v", err),
			ConnectionInfo: u.String(),
		}
	}

	c.Conn = conn
	c.CurrentFeedType = feedType
	return nil
}

func (c *VyperWebsocketClient) Subscribe(feedType FeedType, message interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil {
		return fmt.Errorf("not connected")
	}

	if feedType != c.CurrentFeedType {
		return fmt.Errorf("feed type mismatch")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *VyperWebsocketClient) Unsubscribe(feedType FeedType, message interface{}) error {
	return c.Subscribe(feedType, message)
}

func (c *VyperWebsocketClient) Listen() error {
	if c.Conn == nil {
		return fmt.Errorf("not connected")
	}

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return err
		}

		if c.MessageHandler != nil {
			var rawData map[string]interface{}
			err = json.Unmarshal(message, &rawData)
			if err != nil {
				return err
			}

			convertedData, err := c.convertMessage(rawData)
			if err != nil {
				return err
			}

			c.MessageHandler(convertedData)
		}
	}
}

func (c *VyperWebsocketClient) convertMessage(data map[string]interface{}) (interface{}, error) {
	switch c.CurrentFeedType {
	case WalletEvents:
		return c.convertToChainAction(data)
	case MigrationEvents, TokenEvents:
		return c.convertToTokenPair(data)
	default:
		return nil, fmt.Errorf("unknown feed type: %s", c.CurrentFeedType)
	}
}

func (c *VyperWebsocketClient) convertToChainAction(data map[string]interface{}) (*ChainAction, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var chainAction ChainAction
	err = json.Unmarshal(jsonData, &chainAction)
	if err != nil {
		return nil, err
	}

	return &chainAction, nil
}

func (c *VyperWebsocketClient) convertToTokenPair(data map[string]interface{}) (*TokenPair, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var tokenPair TokenPair
	err = json.Unmarshal(jsonData, &tokenPair)
	if err != nil {
		return nil, err
	}

	return &tokenPair, nil
}

func (c *VyperWebsocketClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil {
		return fmt.Errorf("not connected")
	}

	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	err = c.Conn.Close()
	if err != nil {
		return err
	}

	c.Conn = nil
	c.CurrentFeedType = ""
	return nil
}

func (c *VyperWebsocketClient) Ping() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil {
		return fmt.Errorf("not connected")
	}

	return c.Conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *VyperWebsocketClient) SetMessageHandler(handler MessageHandler) {
	c.MessageHandler = handler
}
