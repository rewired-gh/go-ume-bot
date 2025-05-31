package util

import (
	"bytes"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"

	tg "gopkg.in/telebot.v3"
)

func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ImageToPNGBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetPhotoFileID extracts the file ID from either the current message photo or replied message photo
func GetPhotoFileID(msg *tg.Message) string {
	if msg.Photo != nil {
		return msg.Photo.FileID
	}
	if msg.ReplyTo != nil && msg.ReplyTo.Photo != nil {
		return msg.ReplyTo.Photo.FileID
	}
	return ""
}
