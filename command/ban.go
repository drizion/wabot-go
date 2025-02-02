package command

import (
	"fmt"

	"github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/database/model"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func Ban(msg *events.Message, user *model.BotUser) {
	arg := helpers.GetCmdArgs(msg)[0]

	res, err := client.Wabot.IsOnWhatsApp([]string{arg})
	if err != nil {
		fmt.Println(err)
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Ocorreu um erro ao verificar se o número está no WhatsApp.")
		return
	}

	for _, r := range res {
		if helpers.IsOwner(r.JID) {
			helpers.SendReact(msg, helpers.ForbiddenReaction)
			helpers.Reply(msg, "Você não pode banir o dono do bot.")
			return
		}
		if r.JID.User == client.Wabot.Store.ID.User {
			helpers.SendReact(msg, helpers.ForbiddenReaction)
			helpers.Reply(msg, "Você não pode banir o bot.")
			return
		}
		_, err := client.Wabot.UpdateGroupParticipants(msg.Info.Chat, []types.JID{r.JID}, "remove")
		if err != nil {
			fmt.Println(err)
			helpers.SendReact(msg, helpers.ErrorReaction)
			helpers.Reply(msg, "Ocorreu um erro ao remover o usuário do grupo.")
			return
		} else {
			helpers.ReplyWithMentions(msg, fmt.Sprintf("Membro @%s removido do grupo.", r.JID.User), []string{r.JID.String()})
			helpers.SendReact(msg, helpers.SuccessReaction)
		}
	}
}
