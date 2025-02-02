package command

import (
	"fmt"
	"runtime"
	"time"

	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

const msgFormat = `Pong! 🏓
> dados para nerds 🤓☝️
_WhatsApp:_ %s
_Database:_ %s
_Nº CPUs:_ %d
_Goroutines:_ %d

%d usuários cadastrados`

func Ping(msg *events.Message, user *model.BotUser) {
	startWpp := time.Now()
	helpers.SendReact(msg, helpers.PingReaction)
	durationWpp := time.Since(startWpp)

	startDb := time.Now()
	userCount, _ := helpers.GetUserCount()
	durationDb := time.Since(startDb)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	message := fmt.Sprintf(
		msgFormat,
		durationWpp.String(),
		durationDb.String(),
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		userCount,
	)

	helpers.Reply(msg, message)

}
