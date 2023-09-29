package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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

	re := regexp.MustCompile(`(/?[\w\d@]*\s+?)(.+)`)
	if groups := re.FindStringSubmatch(ctx.Text()); len(groups) >= 3 {
		entity = groups[2]
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
