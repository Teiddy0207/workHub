package constant

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	// Channel related messages
	MESSAGE_CHANNEL   MessageType = "MESSAGE_CHANNEL"
	SUBMIT_FORM       MessageType = "SUBMIT_FORM"
	GET_QUEUE_TICKETS MessageType = "GET_QUEUE_TICKETS"
	CUSTOMER_RATING   MessageType = "CUSTOMER_RATING"
	CURRENT_QUEUE     MessageType = "CURRENT_QUEUE"
	FORM_DATA         MessageType = "FORM_DATA"
)

// BaseMessage is the base structure for all WebSocket messages
type BaseMessage struct {
	Type      MessageType   `json:"type"`
	Request   *FormRequest  `json:"request,omitempty"`
	Response  *FormResponse `json:"response,omitempty"`
	Timestamp time.Time     `json:"timestamp,omitempty"`
}

type FormRequest struct {
	Type     string                 `json:"type"`
	MetaData map[string]interface{} `json:"meta_data"`
}

type FormResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Data     interface{}            `json:"data"`
	Type     string                 `json:"type"`
	MetaData map[string]interface{} `json:"meta_data,omitempty"`
}

// WebSocketClient represents a connected WebSocket client
type WebSocketClient struct {
	ID       string
	UserID   string
	Token    string          // authorization token
	Device   string          // device information from user POST
	Send     chan []byte
}

// ParseMessage parses a JSON message into the BaseMessage struct
func ParseMessage(data []byte) (*BaseMessage, error) {
	var msg BaseMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
