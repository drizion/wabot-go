package command

import (
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Menu(msg *events.Message) {
	helpers.Reply(msg, Registry.GetMenu(msg))
}
