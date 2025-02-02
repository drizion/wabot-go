package helpers

import (
	"context"
	"fmt"

	c "github.com/drizion/wabot-go/client"
	"go.mau.fi/util/random"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"

	"google.golang.org/protobuf/proto"
)

func SendPoll(m *events.Message, text string) {

	optionNames := []string{"Yes", "No"}
	selectableOptionCount := 1
	name := "Casino"

	msgSecret := random.Bytes(32)
	if selectableOptionCount < 0 || selectableOptionCount > len(optionNames) {
		selectableOptionCount = 0
	}
	options := make([]*waProto.PollCreationMessage_Option, len(optionNames))
	for i, option := range optionNames {
		options[i] = &waProto.PollCreationMessage_Option{OptionName: proto.String(option)}
	}
	msg := &waProto.Message{
		PollCreationMessage: &waProto.PollCreationMessage{
			Name:                   proto.String(name),
			Options:                options,
			SelectableOptionsCount: proto.Uint32(uint32(selectableOptionCount)),
			ContextInfo: &waProto.ContextInfo{
				QuotedMessage: m.Message,
			},
		},
		MessageContextInfo: &waProto.MessageContextInfo{
			MessageSecret: msgSecret,
		},
	}

	// var msg = c.Wabot.BuildPollCreation("Casino", []string{"Yes", "No"}, 1)

	_, err := c.Wabot.SendMessage(context.Background(), m.Info.Chat, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
}
