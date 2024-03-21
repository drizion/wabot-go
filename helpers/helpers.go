package helpers

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	c "github.com/drizion/wabot-go/client"
	"go.mau.fi/libsignal/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

const (
	ErrorReaction   string = "❌"
	SuccessReaction string = "✅"
	LoadingReaction string = "⏳"
	PingReaction    string = "🏓"
	LoveReaction    string = "❤️"
	LikeReaction    string = "👍"
	DislikeReaction string = "👎"
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

func WebpWriteExifData(inputData []byte, updateId int64) ([]byte, error) {
	var (
		startingBytes = []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00}
		endingBytes   = []byte{0x16, 0x00, 0x00, 0x00}
		b             bytes.Buffer

		currUpdateId = strconv.FormatInt(updateId, 10)
		currPath     = path.Join("downloads", currUpdateId)
		inputPath    = path.Join(currPath, "input_exif.webm")
		outputPath   = path.Join(currPath, "output_exif.webp")
		exifDataPath = path.Join(currPath, "raw.exif")
	)

	if _, err := b.Write(startingBytes); err != nil {
		return nil, err
	}

	jsonData := map[string]interface{}{
		"sticker-pack-id":        "drizion.dev",
		"sticker-pack-name":      "wabot.net 🤖",
		"sticker-pack-publisher": "@eu_drizion (segue aí rs)",
		"emojis":                 []string{"😀"},
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	jsonLength := (uint32)(len(jsonBytes))
	lenBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuffer, jsonLength)

	if _, err := b.Write(lenBuffer); err != nil {
		return nil, err
	}
	if _, err := b.Write(endingBytes); err != nil {
		return nil, err
	}
	if _, err := b.Write(jsonBytes); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(currPath, os.ModePerm); err != nil {
		return nil, err
	}
	defer os.RemoveAll(currPath)

	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(exifDataPath, b.Bytes(), os.ModePerm); err != nil {
		return nil, err
	}

	cmd := exec.Command("webpmux",
		"-set", "exif",
		exifDataPath, inputPath,
		"-o", outputPath,
	)

	if err := cmd.Run(); err != nil {
		logger.Debug("failed to run webpmux command")

		return nil, err
	}

	return os.ReadFile(outputPath)
}
