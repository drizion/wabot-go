package helpers

import (
	"github.com/drizion/wabot-go/config"
	"go.mau.fi/whatsmeow/types"
)

func IsOwner(sender types.JID) bool {
	for _, num := range config.Bot.OwnerNumbers {
		if sender.User == num {
			return true
		}
	}
	return false
}
