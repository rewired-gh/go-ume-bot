package util

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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

func GetPhotoFileID(msg *tg.Message) string {
	if msg.Photo != nil {
		return msg.Photo.FileID
	}
	if msg.ReplyTo != nil && msg.ReplyTo.Photo != nil {
		return msg.ReplyTo.Photo.FileID
	}
	return ""
}

type ImageFileInfo struct {
	FileID   string
	FileName string
}

func GetImageFileInfo(msg *tg.Message) *ImageFileInfo {
	if msg.Photo != nil {
		return &ImageFileInfo{
			FileID:   msg.Photo.FileID,
			FileName: "",
		}
	}
	if msg.Document != nil && IsAcceptableImageFile(msg.Document.FileName) {
		return &ImageFileInfo{
			FileID:   msg.Document.FileID,
			FileName: msg.Document.FileName,
		}
	}
	if msg.ReplyTo != nil && msg.ReplyTo.Photo != nil {
		return &ImageFileInfo{
			FileID:   msg.ReplyTo.Photo.FileID,
			FileName: "",
		}
	}
	if msg.ReplyTo != nil && msg.ReplyTo.Document != nil && IsAcceptableImageFile(msg.ReplyTo.Document.FileName) {
		return &ImageFileInfo{
			FileID:   msg.ReplyTo.Document.FileID,
			FileName: msg.ReplyTo.Document.FileName,
		}
	}
	return nil
}

// IsAcceptableImageFile checks if the file extension indicates an acceptable image format
func IsAcceptableImageFile(filename string) bool {
	if filename == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(filename))
	acceptableExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".tif"}

	for _, acceptableExt := range acceptableExts {
		if ext == acceptableExt {
			return true
		}
	}

	return false
}
