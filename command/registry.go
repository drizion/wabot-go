package command

import (
	"sync"

	"github.com/drizion/wabot-go/config"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Triggers          []string
	Tags              []string
	MenuTrigger       *string
	Description       string
	Usage             string
	OnlyGroup         bool
	OnlyAdmin         bool
	OnlyPrivate       bool
	OnlyOwner         bool
	HideOnMenu        bool
	AllowUnregistered bool
	Exec              func(msg *events.Message)
}

type CommandRegistry struct {
	commands map[string]Command
	cmdList  []Command
	mu       sync.RWMutex
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

func (r *CommandRegistry) RegisterCommand(command Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cmdList = append(r.cmdList, command)

	for _, trigger := range command.Triggers {
		if _, exists := r.commands[trigger]; exists {
			panic("Command already exists")
		}
		r.commands[trigger] = command
	}

	return nil
}

func (r *CommandRegistry) GetCommand(trigger string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	command, exists := r.commands[trigger]
	return command, exists
}

func (r *CommandRegistry) GetCommandsByTag(tag string) []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var commands []Command
	for _, command := range r.commands {
		for _, t := range command.Tags {
			if t == tag {
				commands = append(commands, command)
				break
			}
		}
	}
	return commands
}

func (r *CommandRegistry) UnregisterCommand(trigger string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.commands, trigger)
}

func (r *CommandRegistry) ListCommands() []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var cmds []Command
	for _, cmd := range r.cmdList {
		if !cmd.HideOnMenu {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

func (r *CommandRegistry) GetMenu() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	menu := helpers.GetGreeting("Fulano") + "\n\n"
	menu += "> *LISTA DE COMANDOS*\n\n"

	for _, command := range r.cmdList {
		if !command.HideOnMenu {
			menu += "`" + config.Prefix + *command.MenuTrigger + "`\n"
			menu += " ↳ _" + command.Description + "_\n\n"
		}
	}
	menu += "Siga o criador do bot no Instagram:\nhttps://instagram.com/eu_drizion"
	return menu
}
