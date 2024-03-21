package handler

import (
	"fmt"

	"github.com/drizion/wabot-go/command"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(m *events.Message) {
	text := GetTextFromMessage(m)
	fmt.Println("Received a message!", text, m.Info)

	if text != "" {
		cmd, exists := command.Registry.GetCommand(text)
		if exists {
			fmt.Println("Command found:", cmd.Triggers)
			cmd.Exec(m)
			return
		}
		fmt.Println("Command not found:", text)
	}

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
