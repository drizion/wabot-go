package command

import (
	"fmt"

	"github.com/drizion/wabot-go/config"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var Registry = NewCommandRegistry()

func SetupCommands() {
	Registry.RegisterCommand(Command{
		MenuTrigger:       proto.String("ping"),
		Triggers:          []string{"ping", "test"},
		Tags:              []string{"test"},
		Exec:              Ping,
		Description:       "Testa a conexão com o bot",
		AllowUnregistered: true,
		Usage:             fmt.Sprintf("Envie `%sping` para testar se o bot está funcionando", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger:       proto.String("menu"),
		Triggers:          []string{"menu", "comandos", "commands"},
		Tags:              []string{"menu", "help", "ajuda", "comandos", "commands"},
		Exec:              Menu,
		Description:       `Exibe a lista de comandos disponíveis.`,
		AllowUnregistered: true,
		Usage:             fmt.Sprintf("Envie `%smenu` para ver a lista de comandos", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger: proto.String("fig"),
		Triggers:    []string{"fig", "figurinhas", "figurinha", "sticker", "stickers"},
		Tags:        []string{"figurinhas"},
		Exec:        func(msg *events.Message) { Fig(msg, nil) },
		Description: "Converte imagem ou video em figurinha.",
		Usage:       fmt.Sprintf("Envie uma imagem ou video com a legenda `%sfig`", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger: proto.String("cfig"),
		Triggers:    []string{"cfig"},
		Tags:        []string{"figurinhas", "figurinhas", "figurinha", "sticker", "stickers"},
		Exec:        func(msg *events.Message) { Fig(msg, []string{"cfig"}) },
		Description: "Converte imagem ou video em figurinha no formato quadradinho (corte automático).",
		Usage:       fmt.Sprintf("Envie uma imagem ou video com a legenda `%scfig`", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger: proto.String("help"),
		Triggers:    []string{"help", "ajuda"},
		Tags:        []string{"help", "ajuda"},
		Exec:        Help,
		Description: "Exibe informações sobre um comando específico.",
		Usage:       fmt.Sprintf("Envie `%shelp <comando>` para ver informações sobre um comando específico", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger: proto.String("ban"),
		Triggers:    []string{"ban"},
		Tags:        []string{"ban"},
		Description: "Banir um usuário.",
		Usage:       fmt.Sprintf("Envie `%sban <usuário>` para banir um usuário", config.Prefix),
		Exec:        Ban,
		OnlyGroup:   true,
		OnlyAdmin:   true,
	})

	fmt.Printf("Registered %d commands\n", len(Registry.cmdList))
}
