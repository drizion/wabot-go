package helpers

import (
	"fmt"

	"github.com/drizion/wabot-go/client"
	"go.mau.fi/whatsmeow/types"
)

func IsAdmin(chat types.JID, sender types.JID) bool {
	res, err := client.Wabot.GetGroupInfo(chat)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, p := range res.Participants {
		if p.JID.User == sender.User {
			return p.IsAdmin || p.IsSuperAdmin
		}
	}
	return false
}
