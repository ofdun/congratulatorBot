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
	bot.Handle("/start", onStart)

	bot.Start()
}

func createMenu(isMailing bool) *tele.ReplyMarkup {
	markup := tele.ReplyMarkup{}
	if isMailing {
		markup.Reply(markup.Row(
			markup.Text("Ура!🎉"),
			markup.Text("Поздравь меня🥳"),
			markup.Text("Отключить рассылку📬"),
		))
		return &markup
	} else {
		markup.Reply(markup.Row(
			markup.Text("Ура!🎉"),
			markup.Text("Поздравь меня🥳"),
			markup.Text("Подключить рассылку📬"),
		))
		return &markup
	}

}

func onStart(c tele.Context) error {
	markup := createMenu(false)
	return c.Send("Привет! Я поздравлю тебя с праздником!", markup)
}

func onText(c tele.Context) error {
	switch c.Text() {
	case "Ура!🎉":
		return c.Send("Ура!🎉")
	case "Поздравь меня🥳":
		{
			postcardPath, err := storage.GetRandomPostcardPath()
			if err != nil {
				return err
			}

			markup := createMenu(false)
			if postcardPath == "" {
				return c.Send("Сегодня нет праздников :(", markup)
			}

			file := tele.FromDisk(postcardPath)
			if download.IsVideo(postcardPath) {
				postcard := &tele.Video{File: file}
				return c.Send(postcard, markup)
			} else {
				postcard := &tele.Photo{File: file}
				return c.Send(postcard, markup)
			}
		}
	case "Подключить рассылку📬":
		{
			return c.Send("Чтобы подключить рассылку напиши /time *time*, например /time 08:30")
		}
	case "Отключить рассылку📬":
		{
			//TODO delete user from db
			return c.Send("Рассылка отключена😈")
		}
	}

	return nil
}
