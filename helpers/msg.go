package helpers

import "go.mau.fi/whatsmeow/types/events"

func GetTextFromMsg(m *events.Message) string {
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
