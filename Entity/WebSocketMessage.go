package Entity

import (
	"github.com/google/uuid"
	"time"
)

type WebsocketMessageRequest struct {
	Id       uuid.UUID   `json:"id"`
	Method   string      `json:"method"`
	DateTime string      `json:"dateTime"`
	Gateway  string      `json:"gateway"`
	Payload  interface{} `json:"payload"`
}
type WebsocketMessageResponse struct {
	Id       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	Status   string    `json:"status"`
	Reason   string    `json:"reason"`
}
