package util

import (
	"path/filepath"

	tg "gopkg.in/telebot.v3"
)

// DownloadFile will download from a given url to a file.
func DownloadFile(directory string, id string, bot *tg.Bot) (filePath string, err error) {
	// Add .jpg extension for downloaded photos to ensure proper format detection
	relativePath := filepath.Join(directory, id+".jpg")
	filePath, err = filepath.Abs(relativePath)
	if err != nil {
		return
	}
	err = bot.Download(&tg.File{FileID: id}, filePath)
	return
}
