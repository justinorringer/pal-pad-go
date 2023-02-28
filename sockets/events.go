package sockets

import (
	"encoding/json"
	"log"

	"github.com/justinorringer/pal-pad-go/db"
)

type eventType string

const (
	sync  eventType = "sync"
	draw  eventType = "draw"
	clear eventType = "clear"
)

type Message struct {
	Type eventType `json:"type"`
	Data []byte    `json:"data"`
}

func processMessage(rc *db.RedisClient, h *Hub, message []byte) (t eventType, err error) {
	m := Message{}
	err = json.Unmarshal(message, &m)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling message: %s", err)
		return
	}

	t = m.Type

	switch m.Type {
	case sync:
		err = db.ProcessSync(rc, m.Data)
	case draw:
		err = db.ProcessDraw(rc, m.Data)
	case clear:
		err = db.ProcessClear(rc, m.Data)
	default:
		// log error
		log.Printf("Unknown event type: %s", m.Type)

		// create error
	}

	return
}
