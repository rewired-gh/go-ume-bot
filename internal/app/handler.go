package app

import (
	"fmt"
	"log/slog"
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
		helpMessage := `基础命令：
/start - 开始和 Ume 聊天啦
/help - 显示此帮助信息
/ping - 测试机器人响应 (回复 pong)

娱乐命令：
/akari - アッカリ〜ン
/hug - 贴贴
/kawaii - 可爱喵
/lu - lu 😭😭
/activate - 赋予上下文
/angry - 😠 (可指定次数)

实用命令：
/n <数字> - 自然数真好玩 (0-10)

图片处理命令：
/burn - 去除盲水印（包括 AI 水印），并破坏图片
/upscale [预设] - 放大图片

图片处理用法：
• 发送图片作为照片或文档，并在标题中添加命令
• 或者回复包含图片的消息并使用命令

可用的放大预设：
• af4 - Anime Fast 4x
• af2 - Anime Fast 2x
• a - Anime Normal 4x
• g - General 4x

不指定预设时默认使用 af4`
		ctx.Reply(helpMessage)
		return nil
	})

	bot.Handle(tg.OnPhoto, func(ctx tg.Context) error {
		return handleMessageCommands(ctx, bot, config)
	})

	bot.Handle(tg.OnDocument, func(ctx tg.Context) error {
		return handleMessageCommands(ctx, bot, config)
	})
}

func downloadImageFromMessage(ctx tg.Context, bot *tg.Bot, config util.Config) (string, error) {
	msg := ctx.Message()
	imageInfo := util.GetImageFileInfo(msg)
	if imageInfo == nil {
		ctx.Reply(util.StickerFromID("CAACAgQAAxkBAAEg1LRkWe1Tk6Vc_mCZ8jqeKN5begPGKgACqwwAAu5XIVKAayOOt2MuRS8E"))
		return "", nil
	}
	return util.DownloadImageFile(config.TmpPath, imageInfo.FileID, imageInfo.FileName, bot)
}

func handleMessageCommands(ctx tg.Context, bot *tg.Bot, config util.Config) error {
	msg := ctx.Message()
	text := util.GetTextOrCaption(msg)

	slog.Debug("Message received", "text", text, "hasPhoto", msg.Photo != nil, "hasDocument", msg.Document != nil)

	command, args := util.ParseCommand(text)

	slog.Debug("Parsed command", "command", command, "args", args)

	switch command {
	case "burn":
		return handleBurnCommand(ctx, bot, config)
	case "upscale":
		return handleUpscaleCommand(ctx, bot, config, args)
	}

	return nil
}

func handleBurnCommand(ctx tg.Context, bot *tg.Bot, config util.Config) error {
	processingMsg, err := ctx.Bot().Reply(ctx.Message(), "正在去除水印，请稍候……")
	if err != nil {
		return err
	}

	filePath, err := downloadImageFromMessage(ctx, bot, config)
	if err != nil {
		bot.Edit(processingMsg, "下载图片失败")
		return err
	}
	if filePath == "" {
		bot.Edit(processingMsg, "没有找到可处理的图片")
		return nil
	}

	img, err := util.LoadImage(filePath)
	if err != nil {
		bot.Edit(processingMsg, "加载图片失败")
		return err
	}
	if err = os.Remove(filePath); err != nil {
		return err
	}
	burnedImage := burner.BurnImage(img)
	imgBytes, err := util.ImageToPNGBytes(burnedImage)
	if err != nil {
		bot.Edit(processingMsg, "处理图片失败")
		return err
	}
	tempFile := fmt.Sprintf("%s/burned_%d.png", config.TmpPath, time.Now().Unix())
	if err := os.WriteFile(tempFile, imgBytes, 0644); err != nil {
		bot.Edit(processingMsg, "保存图片失败")
		return err
	}
	defer os.Remove(tempFile)

	bot.Delete(processingMsg)
	document := &tg.Document{File: tg.FromDisk(tempFile), FileName: "watermark_removed.png"}
	ctx.Reply(document)
	return nil
}

func handleUpscaleCommand(ctx tg.Context, bot *tg.Bot, config util.Config, args string) error {
	slog.Debug("Upscale command started", "args", args)

	processingMsg, err := ctx.Bot().Reply(ctx.Message(), "正在放大图片，请稍候……")
	if err != nil {
		return err
	}

	presetName := util.PresetNameAnimeFast4x
	if args != "" {
		presetName = args
	}
	slog.Debug("Using preset", "preset", presetName)

	filePath, err := downloadImageFromMessage(ctx, bot, config)
	if err != nil {
		slog.Debug("Failed to download file", "error", err)
		bot.Edit(processingMsg, "下载图片失败")
		return err
	}
	if filePath == "" {
		slog.Debug("No image file found")
		bot.Edit(processingMsg, "没有找到可处理的图片")
		return nil
	}
	defer os.Remove(filePath)
	slog.Debug("Downloaded file", "path", filePath)

	slog.Debug("Upscaling image using preset", "preset", presetName)

	resultPath, err := util.UpscaleImage(filePath, presetName, config)
	if err != nil {
		slog.Debug("UpscaleImage failed", "error", err)
		if err == util.ErrInvalidPreset {
			bot.Edit(processingMsg, fmt.Sprintf("无效的预设名称\n\n%s", util.GetPresetsList()))
			return nil
		}
		bot.Edit(processingMsg, "图片放大失败")
		return err
	}
	defer os.Remove(resultPath)
	slog.Debug("Upscaled image saved", "path", resultPath)

	bot.Delete(processingMsg)
	document := &tg.Document{File: tg.FromDisk(resultPath), FileName: fmt.Sprintf("upscaled_%s.png", presetName)}
	ctx.Reply(document)
	slog.Debug("Upscale command completed successfully")

	return nil
}
