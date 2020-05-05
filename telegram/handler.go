package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Tmp pasta, those are already declared in different branch.
type MessageEntities struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Sender struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	Date     int               `json:"date"`
	Entities []MessageEntities `json:"entities"`
	Sender   Sender            `json:"from"`
	Text     string            `json:"text"`
}

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Handler struct {
  service TelegramServiceEntrypoint
}

func NewHandler(s TelegramServiceEntrypoint) Handler {
  return Handler {
    service: s,
  }
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var update Update
  err := decoder.Decode(&update)
  if err != nil {
    fmt.Printf("Failed to decode request %v", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  ctx := context.Background()
  go self.service.ProcessDecodedMessage(ctx, update)
  w.WriteHeader(http.StatusOK)
}
