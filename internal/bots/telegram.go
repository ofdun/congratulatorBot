package bots

import (
	"bufio"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"io"
	"log"
	"os"
	"telegramBot/internal/download"
	"telegramBot/internal/storage"
	"time"
)

type TelegramBot struct {
	*tgbotapi.BotAPI
}

func NewTelegramBot(token string) (TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return TelegramBot{}, err
	}
	return TelegramBot{bot}, err
}

func SetupUpdates(bot *TelegramBot) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func isVideo(path string) bool {
	videoFormats := map[string]bool{
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	}
	suffix := download.GetSuffix(path)
	return videoFormats[suffix]
}

func StartBotPolling(bot *TelegramBot, updates tgbotapi.UpdatesChannel) {
	markup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Поздравь меня🥳"),
			tgbotapi.NewKeyboardButton("Ура!🎉"),
		),
	)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "Поздравь меня🥳":
				{
					waitingMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбираю...")
					waitingMsg.DisableNotification = true

					if _, err := bot.Send(waitingMsg); err != nil {
						panic(err)
					}

					pathToPictureOrVideo, err := chooseRandomPicture()
					if pathToPictureOrVideo != "" {
						start := time.Now().Unix()
						photoOrVideoBytesFile, err := os.Open(pathToPictureOrVideo)
						reader := bufio.NewReader(photoOrVideoBytesFile)

						photoOrVideoBytes, err := io.ReadAll(reader)
						if err != nil {
							panic(err)
						}

						if err != nil {
							panic(err)
						}
						photoOrVideoFileBytes := tgbotapi.FileBytes{
							Name:  "postcard",
							Bytes: photoOrVideoBytes,
						}
						if isVideo(pathToPictureOrVideo) {
							msg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, photoOrVideoFileBytes)
							msg.ReplyMarkup = markup
							_, err = bot.Send(msg)
							println(time.Now().Unix() - start)
							if err != nil {
								panic(err)
							}
						} else {
							// NEW TELEGRAM LIB
							msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photoOrVideoFileBytes)
							msg.ReplyMarkup = markup
							_, err = bot.Send(msg)
							println(time.Now().Unix() - start)
							if err != nil {
								panic(err)
							}
						}

					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сегодня нет празников")
						_, err = bot.Send(msg)
						if err != nil {
							panic(err)
						}
					}
				}
			case "Ура!🎉":
				{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
					_, err := bot.Send(msg)
					if err != nil {
						panic(err)
					}
				}

			}
		}
		//else if update.CallbackQuery != nil {
		//	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		//	if _, err := bot.Send(callback); err != nil {
		//		panic(err)
		//	}
		//
		//	// And finally, send a message containing the data received.
		//	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		//	if _, err := bot.Send(msg); err != nil {
		//		panic(err)
		//	}
		//}
	}
}

func sendKeyboardMessage(bot *tgbotapi.BotAPI, chatID int64, messageText string, markup tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = markup

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func chooseRandomPicture() (string, error) {
	path, err := storage.GetRandomPostcardPath()
	if err != nil {
		return path, err
	}

	return path, nil
}
