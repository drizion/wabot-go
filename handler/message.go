package handler

import (
	"fmt"
	"strings"

	"github.com/drizion/wabot-go/command"
	"github.com/drizion/wabot-go/config"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(m *events.Message) {
	text := GetTextFromMessage(m)

	if m.Info.IsFromMe || text == "" || text[0] != config.Prefix[0] {
		return
	}

	fmt.Println("Received a message!", text)

	text = text[1:]
	trigger := strings.Split(strings.TrimSpace(text), " ")[0]

	cmd, exists := command.Registry.GetCommand(trigger)
	if exists {
		cmd.Exec(m)
		return
	}
	helpers.Reply(m, "Comando não encontrado!\nUse `"+config.Prefix+"menu` para ver a lista de comandos.")
}

func GetTextFromMessage(m *events.Message) string {
	switch {
	case m.Message.GetConversation() != "":
		return m.Message.GetConversation()
	case m.Message.ImageMessage.GetCaption() != "":
		return m.Message.ImageMessage.GetCaption()
	case m.Message.VideoMessage.GetCaption() != "":
		return m.Message.VideoMessage.GetCaption()
	case m.Message.ExtendedTextMessage.GetText() != "":
		return m.Message.ExtendedTextMessage.GetText()
	default:
		return ""
	}
}
