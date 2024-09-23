package vyperclientgo

import "fmt"

type VyperApiError struct {
	Message    string
	StatusCode int
	Response   interface{}
}

func (e *VyperApiError) Error() string {
	return fmt.Sprintf("VyperApiError: %s (Status Code: %d)", e.Message, e.StatusCode)
}

type VyperWebsocketError struct {
	Message        string
	StatusCode     int
	ConnectionInfo interface{}
}

func (e *VyperWebsocketError) Error() string {
	return fmt.Sprintf("VyperWebsocketError: %s (Status Code: %d)", e.Message, e.StatusCode)
}

type AuthenticationError struct {
	VyperApiError
}

type RateLimitError struct {
	VyperApiError
	RetryAfter float64
}

type ServerError struct {
	VyperApiError
}
