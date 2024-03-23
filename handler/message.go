package handler

import (
	"fmt"
	"strings"

	"github.com/drizion/wabot-go/command"
	"github.com/drizion/wabot-go/config"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(m *events.Message) {
	text := helpers.GetTextFromMsg(m)

	if m.Info.IsFromMe || text == "" || text[0] != config.Prefix[0] {
		return
	}

	fmt.Println("Received a message!", text)

	text = text[1:]
	trigger := strings.Split(strings.TrimSpace(text), " ")[0]

	cmd, exists := command.Registry.GetCommand(trigger)
	if exists {
		err := validateConfig(m, &cmd)
		if err != "" {
			helpers.SendReact(m, helpers.ForbiddenReaction)
			helpers.Reply(m, err)
			return
		}
		cmd.Exec(m)
		return
	}
	helpers.Reply(m, "Comando não encontrado!\nUse `"+config.Prefix+"menu` para ver a lista de comandos.")
}

const (
	ErrCmdOnlyGroup   string = "Esse comando só pode ser usado em grupos."
	ErrCmdOnlyOwner   string = "Esse comando só pode ser usado pelo dono do bot."
	ErrCmdOnlyAdmin   string = "Esse comando só pode ser usado por administradores do grupo."
	ErrCmdOnlyPrivate string = "Esse comando só pode ser usado em conversas privadas."
	ErrCmdDisabled    string = "Esse comando está desativado."
)

func validateConfig(m *events.Message, cmd *command.Command) (err string) {
	if cmd.Disabled {
		return ErrCmdDisabled
	}
	if (cmd.OnlyGroup || cmd.OnlyAdmin) && !m.Info.IsGroup {
		return ErrCmdOnlyGroup
	}
	if cmd.OnlyOwner && !helpers.IsOwner(m.Info.Sender) {
		return ErrCmdOnlyOwner
	}
	if cmd.OnlyPrivate && m.Info.IsGroup {
		return ErrCmdOnlyPrivate
	}
	if cmd.OnlyAdmin && !helpers.IsAdmin(m.Info.Chat, m.Info.Sender) {
		return ErrCmdOnlyAdmin
	}

	return ""
}
