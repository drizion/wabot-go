package helpers

import (
	"context"
	"fmt"

	c "github.com/drizion/wabot-go/client"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func Reply(m *events.Message, text string) {
	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      proto.String(m.Info.ID),
				Participant:   proto.String(m.Info.Sender.ToNonAD().String()),
				QuotedMessage: m.Message,
			},
		},
	}

	_, err := c.Wabot.SendMessage(context.Background(), m.Info.Chat, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
}

func ReplyWithMentions(m *events.Message, text string, mentions []string) {
	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      proto.String(m.Info.ID),
				Participant:   proto.String(m.Info.Sender.ToNonAD().String()),
				QuotedMessage: m.Message,
				MentionedJid:  mentions,
			},
		},
	}

	_, err := c.Wabot.SendMessage(context.Background(), m.Info.Chat, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
}

func SendText(from types.JID, text string) {
	var msg = &waProto.Message{
		Conversation: proto.String(text),
	}

	_, err := c.Wabot.SendMessage(context.Background(), from, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
}

func SendTextWithMentions(from types.JID, text string, mentions []string) {
	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				Participant:  proto.String(from.ToNonAD().String()),
				MentionedJid: mentions,
			},
		},
	}

	_, err := c.Wabot.SendMessage(context.Background(), from, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
}
