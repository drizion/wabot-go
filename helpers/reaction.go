package helpers

import (
	"context"
	"fmt"

	c "github.com/drizion/wabot-go/client"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const (
	ErrorReaction     string = "❌"
	ForbiddenReaction string = "🚫"
	SuccessReaction   string = "✅"
	LoadingReaction   string = "⏳"
	ConfigReaction    string = "⚙️"
	PingReaction      string = "🏓"
	LoveReaction      string = "❤️"
	LikeReaction      string = "👍"
	DislikeReaction   string = "👎"
)

func SendReact(m *events.Message, reaction string) whatsmeow.SendResponse {
	r := c.Wabot.BuildReaction(m.Info.Chat, m.Info.Sender, m.Info.ID, reaction)
	resp, err := c.Wabot.SendMessage(context.Background(), m.Info.Chat, r)
	if err != nil {
		fmt.Println("Error sending reaction:", err)
	}
	return resp
}
