package command

import (
	"fmt"

	"github.com/drizion/wabot-go/config"
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
		Usage:             fmt.Sprintf("Envie %sping para testar se o bot está funcionando", config.Prefix),
	})

	Registry.RegisterCommand(Command{
		MenuTrigger:       proto.String("menu"),
		Triggers:          []string{"menu", "help", "ajuda", "comandos", "commands"},
		Tags:              []string{"menu", "help", "ajuda", "comandos", "commands"},
		Exec:              Help,
		Description:       `Exibe a lista de comandos disponíveis.`,
		AllowUnregistered: true,
	})

	Registry.RegisterCommand(Command{
		MenuTrigger: proto.String("fig"),
		Triggers:    []string{"fig", "figurinhas", "figurinha", "sticker", "stickers"},
		Tags:        []string{"figurinhas"},
		Exec:        Fig,
		Description: "Converte imagem ou video em figurinha.",
		Usage:       fmt.Sprintf("Envie uma imagem ou video com a legenda %sfig", config.Prefix),
	})

	fmt.Printf("Registered %d commands\n", len(Registry.cmdList))
}
