package command

import (
	"fmt"

	"github.com/drizion/wabot-go/config"
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func Help(msg *events.Message, user *model.BotUser) {
	arg := helpers.GetCmdArgs(msg)[0]

	if arg == "" {
		helpers.Reply(msg, "Comando não encontrado!\nUse `"+config.Bot.Prefix+"menu` para ver a lista de comandos.")
		return
	}

	cmd, exists := Registry.GetCommand(arg)
	if !exists {
		helpers.Reply(msg, "Comando não encontrado!\nUse `"+config.Bot.Prefix+"menu` para ver a lista de comandos.")
		return
	}

	helpers.Reply(msg, fmt.Sprintf("*Comando `%s%s`*\n%s\n> Uso: %s", config.Bot.Prefix, cmd.Triggers[0], cmd.Description, cmd.Usage))
}
