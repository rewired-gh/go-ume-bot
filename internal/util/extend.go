package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	tg "gopkg.in/telebot.v3"
)

func StickerFromID(id string) *tg.Sticker {
	return &tg.Sticker{
		File: tg.File{
			FileID: id,
		},
	}
}

func RandomPick[T any](arr []T) *T {
	return &arr[rand.Intn(len(arr))]
}

func GetEntity(ctx tg.Context) (entity string) {
	entity = ctx.Sender().FirstName

	if replied := ctx.Message().ReplyTo; replied != nil {
		entity = replied.Sender.FirstName
	}

	re := regexp.MustCompile(`\S+\s+(.*\S)`)
	if groups := re.FindStringSubmatch(ctx.Text()); len(groups) >= 2 {
		entity = groups[1]
	}

	return
}

func SpacifyAfter(str string) string {
	runes := []rune(str)
	lastChar := runes[len(runes)-1]
	if lastChar <= unicode.MaxASCII {
		runes = append(runes, ' ')
	}
	return string(runes)
}

func GetBoundNum(str string) (num int64, err error) {
	num, err = strconv.ParseInt(str, 0, 64)
	if err != nil {
		return
	}
	if num < 0 || num > 10 {
		err = fmt.Errorf("%v is not in valid range", num)
	}
	return
}

func ParseCommand(text string) (command string, args string) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/") {
		return "", ""
	}
	text = text[1:]
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", ""
	}
	command = parts[0]
	if atIndex := strings.Index(command, "@"); atIndex != -1 {
		command = command[:atIndex]
	}
	if len(parts) > 1 {
		args = strings.Join(parts[1:], " ")
	}
	return command, args
}

func GetTextOrCaption(msg *tg.Message) string {
	if msg.Text != "" {
		return msg.Text
	}
	if msg.Caption != "" {
		return msg.Caption
	}
	return ""
}
