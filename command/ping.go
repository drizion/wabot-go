package command

import (
	"time"

	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Ping(msg *events.Message) {
	start := time.Now()
	helpers.SendReact(msg, helpers.PingReaction)
	duration := time.Since(start)

	helpers.Reply(msg, "Pong! 🏓\n_Latência no WhatsApp:_ "+duration.String())
}
