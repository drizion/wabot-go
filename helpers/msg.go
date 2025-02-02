package helpers

import (
	"regexp"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

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

func GetCmdArgs(m *events.Message) []string {
	argsString := strings.Join(strings.Fields(GetTextFromMsg(m)[1:]), " ")
	re := regexp.MustCompile(`\s+`)
	args := re.Split(argsString, -1)[1:]

	if len(args) == 0 {
		args = append(args, "")
	}

	return args
}
