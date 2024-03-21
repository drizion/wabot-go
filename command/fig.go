package command

import (
	"context"
	"fmt"
	"os"
	"time"

	c "github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/helpers"
	"github.com/h2non/bimg"
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

	stickerBytes, err = helpers.WebpWriteExifData(newImage, msg.Info.Timestamp.Unix())
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
