package bots

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"os"
	"telegramBot/internal/storage"
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

func StartBotPolling(bot *TelegramBot, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			// update.Message.Text == msg text from user
			// update.Message.Chat.ID == chat id

			switch update.Message.Text {
			case "Поздравь меня🥳":
				{

					pathToPicture, err := chooseRandomPicture()
					if pathToPicture != "" {
						photoBytes, err := os.ReadFile(pathToPicture)
						if err != nil {
							panic(err)
						}
						photoFileBytes := tgbotapi.FileBytes{
							Name:  "postcard",
							Bytes: photoBytes,
						}
						msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photoFileBytes)
						_, err = bot.Send(msg)
						if err != nil {
							panic(err)
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
	}
}

func chooseRandomPicture() (string, error) {
	path, err := storage.GetRandomPostcardPath()
	if err != nil {
		return path, err
	}

	return path, nil
}
