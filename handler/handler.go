package handler

import (
	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go MessageHandler(v)
	}
}
