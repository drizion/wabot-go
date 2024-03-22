package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	c "github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/helpers"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func Fig(msg *events.Message, args []string) {
	go helpers.SendReact(msg, helpers.LoadingReaction)

	var isVideoFig bool = false

	now := time.Now()
	cropSquare := false

	for _, arg := range args {
		if arg == "cfig" {
			cropSquare = true
		}
	}

	var imageMsg *waProto.ImageMessage
	var videoMsg *waProto.VideoMessage
	var uploadedSticker whatsmeow.UploadResponse
	var stickerBytes []byte

	imageMsg, err := helpers.GetImageMessage(msg)
	if err != nil {
		videoMsg, err = helpers.GetVideoMessage(msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			helpers.SendReact(msg, helpers.ErrorReaction)
			helpers.Reply(msg, "Você precisa enviar uma imagem ou um vídeo para que eu possa transformá-lo em uma figurinha.")
			return
		}
		isVideoFig = true
	}

	if isVideoFig {
		stickerBytes, err = c.Wabot.Download(videoMsg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			helpers.SendReact(msg, helpers.ErrorReaction)
			helpers.Reply(msg, "Ocorreu um erro ao baixar o vídeo... por favor, tente novamente.")
			return
		}
	} else {
		stickerBytes, err = c.Wabot.Download(imageMsg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			helpers.SendReact(msg, helpers.ErrorReaction)
			helpers.Reply(msg, "Ocorreu um erro ao baixar a imagem... por favor, tente novamente.")
			return
		}
	}

	inputPath := fmt.Sprintf("temp/%d.gif", now.Unix())
	outputPath := fmt.Sprintf("temp/%d.webp", now.Unix())

	_ = os.WriteFile(inputPath, stickerBytes, 0644)

	if cropSquare {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libwebp", "-preset", "default", "-loop", "0", "-an", "-qscale:v", "20", "-t", "00:00:10", "-compression_level", "100", "-vf", "crop=in_w:in_w:0:(in_h-in_w)/2,scale=320:320,fps=20", outputPath)
		if err := cmd.Run(); err != nil {
			fmt.Println("Erro ao converter o arquivo:", err)
			helpers.Reply(msg, "Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
			_ = os.Remove(inputPath)
			_ = os.Remove(outputPath)
			return
		}
	} else {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libwebp", "-preset", "default", "-loop", "0", "-an", "-qscale:v", "20", "-t", "00:00:10", "-compression_level", "100", "-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=20, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse", outputPath)
		if err := cmd.Run(); err != nil {
			fmt.Println("Erro ao converter o arquivo:", err)
			helpers.Reply(msg, "Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
			_ = os.Remove(inputPath)
			_ = os.Remove(outputPath)
			return
		}
	}

	fmt.Println("Conversão concluída com sucesso!")

	stickerBytes, err = os.ReadFile(outputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
		return
	}

	_ = os.Remove(inputPath)
	_ = os.Remove(outputPath)

	// if stickerBytes length > 1mb, return error
	if len(stickerBytes) > 1024*1024 {
		fmt.Println("Sticker has a lot of data which cannot be handled by WhatsApp")
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "A figurinha é muito pesada, tente diminuir o tamanho ou a resolução e envie novamente.")
		return
	}

	stickerBytes, err = helpers.WebpWriteExifData(stickerBytes, now.Unix())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Ocorreu um erro ao adicionar os metadados... por favor, tente novamente.")
		return
	}

	uploadedSticker, err = c.Wabot.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Ocorreu um erro ao fazer o upload do video... por favor, tente novamente.")
		return
	}

	fmt.Println("Sticker uploaded:", uploadedSticker)
	msgToSend := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploadedSticker.URL),
			DirectPath:    proto.String(uploadedSticker.DirectPath),
			MediaKey:      uploadedSticker.MediaKey,
			IsAnimated:    proto.Bool(bool(isVideoFig)),
			IsAvatar:      proto.Bool(false),
			Mimetype:      proto.String("image/webp"),
			FileEncSha256: uploadedSticker.FileEncSHA256,
			FileSha256:    uploadedSticker.FileSHA256,
			FileLength:    proto.Uint64(uploadedSticker.FileLength),
			StickerSentTs: proto.Int64(time.Now().Unix()),
		},
	}

	go helpers.SendReact(msg, helpers.SuccessReaction)

	_, err = c.Wabot.SendMessage(context.Background(), msg.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		helpers.SendReact(msg, helpers.ErrorReaction)
		helpers.Reply(msg, "Ocorreu um erro ao enviar a figurinha... por favor, tente novamente.")
	}
}
