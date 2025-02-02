package helpers

import (
	db "github.com/drizion/wabot-go/database"
	"github.com/drizion/wabot-go/database/model"
)

func GetUserCount() (int64, error) {

	var count int64

	res := db.DB.Model(&model.BotUser{}).Count(&count)

	if res.Error != nil {
		return 0, res.Error
	}

	return count, nil
}
