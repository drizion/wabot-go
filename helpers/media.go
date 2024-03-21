package helpers

import (
	"errors"
	"fmt"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

const (
	Image   = "image"
	Video   = "video"
	Audio   = "audio"
	Sticker = "sticker"
)

type MessageType string

const (
	ImageMessageType   MessageType = "image"
	VideoMessageType   MessageType = "video"
	AudioMessageType   MessageType = "audio"
	StickerMessageType MessageType = "sticker"
)

func GetMessage(msg *events.Message, msgType MessageType) (message interface{}, err error) {
	switch msgType {
	case "image":
		message, err = GetImageMessage(msg)
	case "video":
		message, err = GetVideoMessage(msg)
	case "audio":
		message, err = GetAudioMessage(msg)
	case "sticker":
		message, err = GetStickerMessage(msg)
	default:
		return nil, errors.New("unsupported message type")
	}

	if message == nil || err != nil {
		return nil, fmt.Errorf("no valid %s found", msgType)
	}

	return message, nil
}

func GetImageMessage(msg *events.Message) (imageMsg *waProto.ImageMessage, err error) {
	msg.SourceWebMsg.GetMediaData()
	if msg.Message.ImageMessage != nil {
		imageMsg = msg.Message.ImageMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage != nil {
		imageMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage
	} else {
		err = fmt.Errorf("no valid image found")
		return nil, err
	}
	return imageMsg, nil
}

func GetVideoMessage(msg *events.Message) (videoMsg *waProto.VideoMessage, err error) {
	if msg.Message.VideoMessage != nil {
		videoMsg = msg.Message.VideoMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage != nil {
		videoMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage
	} else {
		err = fmt.Errorf("no valid video found")
		return nil, err
	}
	return videoMsg, nil
}

func GetAudioMessage(msg *events.Message) (audioMsg *waProto.VideoMessage, err error) {
	if msg.Message.VideoMessage != nil {
		audioMsg = msg.Message.VideoMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage != nil {
		audioMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage
	} else {
		err = fmt.Errorf("no valid video found")
		return nil, err
	}
	return audioMsg, nil
}

func GetStickerMessage(msg *events.Message) (stickerMsg *waProto.StickerMessage, err error) {
	if msg.Message.StickerMessage != nil {
		stickerMsg = msg.Message.StickerMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.StickerMessage != nil {
		stickerMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.StickerMessage
	} else {
		err = fmt.Errorf("no valid video found")
		return nil, err
	}
	return stickerMsg, nil
}
