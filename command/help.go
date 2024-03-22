package command

import (
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

func Help(msg *events.Message) {
	fmt.Println("Command HELP executed")
}
