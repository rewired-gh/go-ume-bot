package app

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/rewired-gh/go-ume-bot/internal/util"
	"github.com/rewired-gh/water-burner/pkg/burner"
	tg "gopkg.in/telebot.v3"
)

func HandleCommands(bot *tg.Bot, config util.Config) {
	bot.Handle("/start", func(ctx tg.Context) error {
		ctx.Reply("Hi~")
		return nil
	})

	bot.Handle("/activate", func(ctx tg.Context) error {
		ctx.Reply("现在我获得自由说话的神丹了")
		return nil
	})

	bot.Handle("/akari", func(ctx tg.Context) error {
		ctx.Reply(util.StickerFromID("CAACAgIAAxkBAAEZ-Zpjc3WsWVrRzMzls9EdsB1sXWte1AACMisAAuCjggerIQY25DgHESsE"))
		return nil
	})

	bot.Handle("/hug", func(ctx tg.Context) error {
		stickers := []string{
			"CAACAgQAAxkBAAEZ-gJjc4PDazUpm3xSxvDCLjJ_HxoSwgAC6wADo30xFaajclVMkzh7KwQ",
			"CAACAgQAAxkBAAEZ-gtjc4SlPDR-PlC8zJrCXFDSjjI6jQACUAIAAqN9MRV7rBhq5jvKkCsE",
			"CAACAgUAAxkBAAEZ-g1jc4TQMXGbT_7c9qJFBCIeIa4vRQACWAUAAp0xYFY-jmPr-ng3iysE",
			"CAACAgQAAxkBAAEZ-g9jc4TvudlCbLrD5h_cXDo1KnyOzQAC3QoAAkkoUFPEHy70_HWcLysE",
		}
		ctx.Reply(fmt.Sprintf("%s贴贴！", util.SpacifyAfter(util.GetEntity(ctx))))
		ctx.Reply(util.StickerFromID(*util.RandomPick(stickers)))
		return nil
	})

	bot.Handle("/kawaii", func(ctx tg.Context) error {
		ctx.Reply(fmt.Sprintf("%s可爱喵！", util.SpacifyAfter(util.GetEntity(ctx))))
		return nil
	})

	bot.Handle("/lu", func(ctx tg.Context) error {
		quotes := []string{
			"你到底是谁😭",
			"你在哪里😭",
			"你带我走吧😭",
			"你给我出来😭",
		}
		quote := *util.RandomPick(quotes)
		ctx.Reply(fmt.Sprintf("%s%s", util.SpacifyAfter(util.GetEntity(ctx)), quote))
		return nil
	})

	bot.Handle("/angry", func(ctx tg.Context) error {
		init := "😠😠😠"
		states := []string{"😡😠😠", "😠😡😠", "😠😠😡"}
		num, err := util.GetBoundNum(util.GetEntity(ctx))
		if err != nil {
			num = 6
		}
		msg, err := bot.Reply(ctx.Message(), init)
		if err != nil {
			return err
		}
		for i := int64(0); i < num; i++ {
			time.Sleep(2 * time.Second)
			bot.Edit(msg, states[i%int64(len(states))])
		}
		return nil
	})

	bot.Handle("/ping", func(ctx tg.Context) error {
		ctx.Reply("pong")
		return nil
	})

	bot.Handle("/n", func(ctx tg.Context) error {
		num, err := util.GetBoundNum(ctx.Message().Payload)
		if err != nil {
			ctx.Reply(util.StickerFromID("CAACAgQAAxkBAAEg1LRkWe1Tk6Vc_mCZ8jqeKN5begPGKgACqwwAAu5XIVKAayOOt2MuRS8E"))
			return err
		}
		ctx.Reply(util.GetNaturalSet(num))
		return nil
	})

	bot.Handle("/burn", func(ctx tg.Context) error {
		return handleBurnCommand(ctx, bot, config)
	})

	bot.Handle("/upscale", func(ctx tg.Context) error {
		return handleUpscaleCommand(ctx, bot, config, ctx.Message().Payload)
	})

	bot.Handle("/help", func(ctx tg.Context) error {
		helpMessage := `/start - 开始和 Ume 聊天啦
/help - 显示帮助信息
/akari - アッカリ〜ン
/hug - 贴贴
/kawaii - 可爱喵
/lu - lu 😭😭
/activate - 赋予上下文
/angry - 😠
/n - 自然数真好玩
/ping - pong
/burn - 去除盲水印（包括 AI 水印），并破坏图片
/upscale - 放大图片 (可选预设名称)`
		ctx.Reply(helpMessage)
		return nil
	})

	// Handle both photo and text messages with a unified handler
	bot.Handle(tg.OnPhoto, func(ctx tg.Context) error {
		return handleMessageCommands(ctx, bot, config)
	})

	bot.Handle(tg.OnText, func(ctx tg.Context) error {
		return handleMessageCommands(ctx, bot, config)
	})
}

// handleMessageCommands processes commands in both photo captions and text messages
func handleMessageCommands(ctx tg.Context, bot *tg.Bot, config util.Config) error {
	msg := ctx.Message()
	text := util.GetTextOrCaption(msg)

	command, args := util.ParseCommand(text)

	switch command {
	case "burn":
		return handleBurnCommand(ctx, bot, config)
	case "upscale":
		return handleUpscaleCommand(ctx, bot, config, args)
	}

	return nil
}

// handleBurnCommand handles burn command for any message type
func handleBurnCommand(ctx tg.Context, bot *tg.Bot, config util.Config) error {
	msg := ctx.Message()
	fileID := util.GetPhotoFileID(msg)
	if fileID == "" {
		ctx.Reply(util.StickerFromID("CAACAgQAAxkBAAEg1LRkWe1Tk6Vc_mCZ8jqeKN5begPGKgACqwwAAu5XIVKAayOOt2MuRS8E"))
		return nil
	}
	filePath, err := util.DownloadFile(config.TmpPath, fileID, bot)
	if err != nil {
		return err
	}
	img, err := util.LoadImage(filePath)
	if err != nil {
		return err
	}
	if err = os.Remove(filePath); err != nil {
		return err
	}
	burnedImage := burner.BurnImage(img)
	imgBytes, err := util.ImageToPNGBytes(burnedImage)
	if err != nil {
		return err
	}
	photo := &tg.Photo{File: tg.FromReader(bytes.NewReader(imgBytes))}
	ctx.Reply(photo)
	return nil
}

// handleUpscaleCommand handles upscale command for any message type
func handleUpscaleCommand(ctx tg.Context, bot *tg.Bot, config util.Config, args string) error {
	msg := ctx.Message()
	presetName := util.PresetNameAnimeFast4x
	if args != "" {
		presetName = args
	}

	fileID := util.GetPhotoFileID(msg)
	if fileID == "" {
		ctx.Reply(util.StickerFromID("CAACAgQAAxkBAAEg1LRkWe1Tk6Vc_mCZ8jqeKN5begPGKgACqwwAAu5XIVKAayOOt2MuRS8E"))
		return nil
	}

	filePath, err := util.DownloadFile(config.TmpPath, fileID, bot)
	if err != nil {
		return err
	}
	defer os.Remove(filePath)

	resultPath, err := util.UpscaleImage(filePath, presetName, config)
	if err != nil {
		if err == util.ErrInvalidPreset {
			ctx.Reply(util.GetPresetsList())
			return nil
		}
		return err
	}
	defer os.Remove(resultPath)

	photo := &tg.Photo{File: tg.FromDisk(resultPath)}
	ctx.Reply(photo)

	return nil
}
