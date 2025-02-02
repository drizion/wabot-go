package command

import (
	"fmt"

	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Money(msg *events.Message, user *model.BotUser) {
	helpers.Reply(msg, fmt.Sprintf("> Saldo 💰\n```%d botcoins```", user.WaCoins))
}
