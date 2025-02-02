package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/command"
	"github.com/drizion/wabot-go/config"
	db "github.com/drizion/wabot-go/database"
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(m *events.Message) {
	text := helpers.GetTextFromMsg(m)

	pollUpdate := m.Message.GetPollUpdateMessage()

	if pollUpdate != nil {
		go PollUpdateMessageHandler(m)
		return
	}

	if text == "" || text[0] != config.Bot.Prefix[0] {
		return
	}

	if m.Info.IsGroup {
		var group model.Group

		result := db.DB.First(&group, "id = ?", m.Info.Chat)

		if result.Error != nil {
			fmt.Println("Error creating group", result.Error)
		} else {
			fmt.Println("Group created, rows affected", result.RowsAffected)
		}

	}

	fmt.Println("Received a message!", text)

	text = text[1:]
	trigger := strings.Split(strings.TrimSpace(text), " ")[0]

	user, _ := helpers.GetUser(m.Info.Sender)

	cmd, exists := command.Registry.GetCommand(trigger)
	if exists {
		err := validateConfig(m, &cmd, &user)
		if err != "" {
			helpers.SendReact(m, helpers.ForbiddenReaction)
			helpers.Reply(m, err)
			return
		}
		cmd.Exec(m, &user)
		return
	}
	helpers.Reply(m, "Comando não encontrado!\nUse `"+config.Bot.Prefix+"menu` para ver a lista de comandos.")
}

const (
	ErrCmdOnlyGroup   string = "Esse comando só pode ser usado em grupos."
	ErrCmdOnlyOwner   string = "Esse comando só pode ser usado pelo dono do bot."
	ErrCmdOnlyAdmin   string = "Esse comando só pode ser usado por administradores do grupo."
	ErrCmdOnlyPrivate string = "Esse comando só pode ser usado em conversas privadas."
	ErrCmdDisabled    string = "Esse comando está desativado."
)

func validateConfig(m *events.Message, cmd *command.Command, user *model.BotUser) (err string) {
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
	if !cmd.AllowUnregistered && user.ID == "" {
		return "> ⓘ Você não está registrado no sistema. Use `" + config.Bot.Prefix + "register` para se registrar."
	}

	return ""
}

func PollUpdateMessageHandler(m *events.Message) {
	poll, err := client.Wabot.DecryptPollVote(m)

	if err != nil {
		fmt.Printf("Error decrypting poll: %v", err)
		return
	}

	if poll != nil && len(poll.SelectedOptions) > 0 {

		jsonM, err := json.MarshalIndent(m, "", "  ")

		if err != nil {
			fmt.Printf("Error marshalling poll: %v", err)
			return
		}

		fmt.Println("Poll update received!", string(jsonM))
		helpers.SendText(m.Info.Chat, fmt.Sprintf("enquete votada no id %s", *m.Message.PollUpdateMessage.PollCreationMessageKey.Id))
		// helpers.SendText(m.Info.Chat, string(jsonM))

		if string(poll.SelectedOptions[0]) == string(command.BetAcceptedHash) {
			helpers.SendTextWithMentions(m.Info.Chat, fmt.Sprintf("🎰 @%s aceitou a aposta!\njoguei os dados no pv 🎲", m.Info.Sender.User), []string{m.Info.Sender.User + "@" + m.Info.Sender.Server})

			dice1 := helpers.GetRandomD6()
			dice2 := helpers.GetRandomD6()

			time.Sleep(1 * time.Second)

			helpers.SendText(m.Info.Sender, fmt.Sprintf("🎲 Dado 1: %s", dice1.Emoji))
			helpers.SendText(m.Info.Sender, fmt.Sprintf("🎲 Dado 2: %s", dice2.Emoji))

			helpers.SendText(m.Info.Sender, fmt.Sprintf("> 🎲 Total: %d", dice1.Value+dice2.Value))

			time.Sleep(2 * time.Second)

			if dice1.Value+dice2.Value < 12 {
				pollOptions := []string{"Arriscar", "Tô de boa"}

				poll := client.Wabot.BuildPollCreation(
					"🎲 Deseja jogar mais um dado? se o valor total passar de 12, você perde...",
					pollOptions,
					1,
				)

				hashes := command.HashPollOptions(pollOptions)

				fmt.Println("Poll options hash", hashes)

				jsonObj, err := json.MarshalIndent(hashes, "", "  ")
				if err != nil {
					fmt.Printf("Error marshalling poll: %v", err)
					return
				}

				fmt.Println("Poll creation", string(jsonObj))

				_, err = client.Wabot.SendMessage(context.Background(), m.Info.Sender, poll)

				if err != nil {
					fmt.Printf("Error sending message: %v", err)
				}
			}
			return
		}

		if string(poll.SelectedOptions[0]) == string(command.BetCanceledHash) {
			helpers.Reply(m, "Aposta cancelada! 🎰")
			return
		}

		jsonObj, err := json.MarshalIndent(poll, "", "  ")

		if err != nil {
			fmt.Printf("Error marshalling poll: %v", err)
			return
		}

		fmt.Println("Poll update received!", string(jsonObj))

	}
}
