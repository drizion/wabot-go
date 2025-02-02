package command

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	db "github.com/drizion/wabot-go/database"
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"github.com/google/uuid"
	"go.mau.fi/whatsmeow/types/events"
	"gorm.io/gorm"
)

func Bonus(msg *events.Message, user *model.BotUser) {
	var dailyBonus model.DailyBonu
	now := time.Now()

	res := db.DB.Last(&dailyBonus, model.DailyBonu{
		BotUserID: user.ID,
	})

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println(res.Error)
		helpers.Reply(msg, "Erro ao buscar bônus diário")
		return
	}

	// if the user has already received the bonus today, return an error message
	if dailyBonus.CreatedAt.Format("02/01/2006") == now.Format("02/01/2006") {
		helpers.Reply(msg, "Você já pegou seu bônus diário hoje")
		return
	} else {
		min := 0
		max := 10

		random := int32(rand.Intn(max-min) + min)

		id, err := uuid.NewV7()
		if err != nil {
			helpers.Reply(msg, "Erro ao gerar ID")
			return
		}

		dailyBonus.BotUserID = user.ID
		dailyBonus.ID = id.String()
		dailyBonus.CreatedAt = now
		db.DB.Create(&dailyBonus)

		user.WaCoins += random
		go db.DB.Save(&user)

		go helpers.Reply(msg, fmt.Sprintf("Você pegou seu bônus diário de %d botcoins! ", random))
		return
	}

}
