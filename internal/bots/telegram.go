package bots

import (
	tele "gopkg.in/telebot.v3"
	"telegramBot/internal/download"
	"telegramBot/internal/storage"
	"time"
)

type TelegramBot struct {
	*tele.Bot
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	return &TelegramBot{bot}, nil
}

func StartTelegramBot(bot *TelegramBot) {
	bot.Handle(tele.OnText, onText)

	bot.Start()
}

func onText(c tele.Context) error {
	switch c.Text() {
	case "Ура!🎉":
		return c.Send("Ура!🎉")
	case "Поздравь меня🥳":
		postcardPath, err := storage.GetRandomPostcardPath()
		if err != nil {
			return err
		}

		if postcardPath == "" {
			return c.Send("Сегодня нет праздников :(")
		}

		file := tele.FromDisk(postcardPath)
		if download.IsVideo(postcardPath) {
			postcard := &tele.Video{File: file}
			return c.Send(postcard)
		} else {
			postcard := &tele.Photo{File: file}
			return c.Send(postcard)
		}
	}
	return nil
}
