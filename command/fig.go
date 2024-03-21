package command

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
	"github.com/drizion/wabot-go/helpers"
	"github.com/h2non/bimg"
	"go.mau.fi/libsignal/logger"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func Fig(msg *events.Message) {
	fmt.Println("Command FIG executed")

	var imageMsg *waProto.ImageMessage

	if msg.Message.ImageMessage != nil {
		imageMsg = msg.Message.ImageMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage != nil {
		imageMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage
	} else {
		// Nenhum dos campos contém uma imagem válida
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Você precisa enviar uma imagem para que eu possa transformá-la em um sticker.")
		return
	}

	stickerBytes, err := c.Wabot.Download(imageMsg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	newImage, err := bimg.NewImage(stickerBytes).Convert(bimg.WEBP)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	stickerBytes, err = WebpWriteExifData(newImage, msg.Info.Timestamp.Unix())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	uploadedSticker, err := c.Wabot.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("Sticker uploaded:", uploadedSticker)

	msgToSend := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploadedSticker.URL),
			DirectPath:    proto.String(uploadedSticker.DirectPath),
			MediaKey:      uploadedSticker.MediaKey,
			IsAnimated:    proto.Bool(bool(false)),
			IsAvatar:      proto.Bool(false),
			Height:        proto.Uint32(uint32(*imageMsg.Height)),
			Width:         proto.Uint32(uint32(*imageMsg.Width)),
			Mimetype:      proto.String("image/webp"),
			FileEncSha256: uploadedSticker.FileEncSHA256,
			FileSha256:    uploadedSticker.FileSHA256,
			FileLength:    proto.Uint64(uploadedSticker.FileLength),
			StickerSentTs: proto.Int64(time.Now().Unix()),
		},
	}
	res, err := c.Wabot.SendMessage(context.Background(), msg.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("Message sent:", res)
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
