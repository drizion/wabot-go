package helpers

import (
	"errors"

	db "github.com/drizion/wabot-go/database"
	"github.com/drizion/wabot-go/database/model"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

func GetUser(sender types.JID) (model.BotUser, error) {
	var user model.BotUser

	id := sender.User + "@" + sender.Server

	result := db.DB.First(&user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, result.Error
	}

	return user, nil
}
