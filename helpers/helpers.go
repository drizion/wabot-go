package helpers

import (
	"context"
	"fmt"
	"time"

	c "github.com/drizion/wabot-go/client"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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

func GetGreeting(sender string) string {
	// Configurar o locale para pt-br
	// dayjs.locale('pt-br');
	// Obter a data e hora atual
	now := time.Now().UTC().Add(-3 * time.Hour)
	// Verificar o período do dia e retornar a saudação apropriada
	hour := now.Hour()
	if hour >= 4 && hour < 12 {
		return "Bom dia " + sender + ", dormiu bem?"
	} else if hour >= 12 && hour < 18 {
		if hour == 12 {
			return "Boa tarde " + sender + ", já almoçou?"
		}
		return "Boa tarde " + sender + ", como vai?"
	} else {
		return "Boa noite " + sender + ", tudo bem?"
	}
}

func SendReact(m *events.Message, reaction string) whatsmeow.SendResponse {
	r := c.Wabot.BuildReaction(m.Info.Chat, m.Info.Sender, m.Info.ID, reaction)
	resp, err := c.Wabot.SendMessage(context.Background(), m.Info.Chat, r)
	if err != nil {
		fmt.Println("Error sending reaction:", err)
	}
	return resp
}

const (
	ErrorReaction   string = "❌"
	SuccessReaction string = "✅"
	LoadingReaction string = "⏳"
	PingReaction    string = "🏓"
	LoveReaction    string = "❤️"
	LikeReaction    string = "👍"
	DislikeReaction string = "👎"
)
