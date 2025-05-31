package util

import (
	"path/filepath"

	tg "gopkg.in/telebot.v3"
)

func DownloadImageFile(directory string, id string, filename string, bot *tg.Bot) (filePath string, err error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}

	relativePath := filepath.Join(directory, id+ext)
	filePath, err = filepath.Abs(relativePath)
	if err != nil {
		return
	}
	err = bot.Download(&tg.File{FileID: id}, filePath)
	return
}
