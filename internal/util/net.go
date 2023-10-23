package util

import (
	"path/filepath"

	tg "gopkg.in/telebot.v3"
)

// DownloadFile will download from a given url to a file.
func DownloadFile(directory string, id string, bot *tg.Bot) (filePath string, err error) {
	filePath = filepath.Join(directory, id)
	err = bot.Download(&tg.File{FileID: id}, filePath)
	return
}
