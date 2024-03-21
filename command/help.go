package command

import (
	"fmt"

	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Help(msg *events.Message) {
	fmt.Println("Command HELP executed")

	helpers.Reply(msg, Registry.GetMenu())
}
