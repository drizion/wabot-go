package command

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	c "github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/config"
	db "github.com/drizion/wabot-go/database"
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	BetAcceptedHash = []byte{
		61, 248, 100, 74, 251, 23, 66, 221,
		90, 5, 87, 208, 123, 73, 136, 179,
		192, 105, 137, 53, 183, 67, 178, 228,
		179, 186, 56, 48, 184, 117, 127, 140,
	}
	BetCanceledHash = []byte{
		2, 131, 171, 101, 18, 121, 20, 29,
		143, 53, 82, 179, 44, 31, 70, 25,
		217, 82, 148, 19, 35, 152, 107, 46,
		227, 224, 98, 66, 168, 28, 216, 132,
	}
)

func Casino(msg *events.Message, user *model.BotUser) {
	arg := helpers.GetCmdArgs(msg)[0]

	if arg == "" {
		helpers.Reply(msg, fmt.Sprintf("Você precisa especificar o valor da aposta!\nUse `%scasino <valor>` para apostar.", config.Bot.Prefix))
		return
	}

	amount, err := strconv.Atoi(arg)

	if err != nil {
		helpers.Reply(msg, "Valor inválido!")
		return
	}

	if amount > int(user.WaCoins) {
		helpers.Reply(msg, fmt.Sprintf("Você não tem saldo suficiente!\nSeu saldo atual é de %d botcoins.", user.WaCoins))
		return
	}

	currency := "botcoin"

	if amount > 1 {
		currency += "s"
	}

	if amount < 1 {
		helpers.Reply(msg, "Você deve apostar no mínimo 1 botcoin.")
		return
	}

	pollOptions := []string{"Aceitar aposta", "Cancelar aposta"}

	poll := c.Wabot.BuildPollCreation(
		fmt.Sprintf("🎰 Apostar %d %s 🎰", amount, currency),
		pollOptions,
		1,
	)

	hashes := HashPollOptions(pollOptions)

	fmt.Println("Poll options hash", hashes)

	jsonObj, err := json.MarshalIndent(hashes, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling poll: %v", err)
		return
	}

	fmt.Println("Poll creation", string(jsonObj))

	pollMessage, err := c.Wabot.SendMessage(context.Background(), msg.Info.Chat, poll)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	now := time.Now()

	db.DB.Create(&model.CasinoSession{
		ID:        pollMessage.ID,
		BotUserID: user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		Amount:    int32(amount),
		GroupID:   msg.Info.Chat.User + "@" + msg.Info.Chat.Server,
		Status:    "pending",
	})

	helpers.Reply(msg, fmt.Sprintf("enquete criada no ID %s", pollMessage.ID))
}

func HashPollOptions(optionNames []string) [][]byte {
	optionHashes := make([][]byte, len(optionNames))
	for i, option := range optionNames {
		optionHash := sha256.Sum256([]byte(option))
		optionHashes[i] = optionHash[:]
	}
	return optionHashes
}
