package command

import (
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Menu(msg *events.Message, user *model.BotUser) {
	helpers.Reply(msg, Registry.GetMenu(msg))
}
