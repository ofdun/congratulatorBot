package bots

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strconv"
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
	bot.Handle("/time", onTime)

	bot.Start()
}

func timeIsValid(time string) (bool, string) {
	if len(time) < 4 || len(time) > 5 {
		return false, ""
	} else if len(time) == 4 {
		time = "0" + time
	}
	if _, err := strconv.Atoi(string(time[2])); err == nil {
		return false, ""
	}

	hoursString := time[0:2]
	minutesString := time[3:5]

	hours, err := strconv.Atoi(hoursString)
	if err != nil {
		return false, ""
	}
	if hours > 23 || hours < 0 {
		return false, ""
	}

	minutes, err := strconv.Atoi(minutesString)
	if err != nil {
		return false, ""
	}
	if minutes > 59 || minutes < 0 {
		return false, ""
	}

	return true, hoursString + ":" + minutesString
}

func onTime(c tele.Context) error {
	mailingTime := c.Message().Payload
	if valid, formattedTime := timeIsValid(mailingTime); valid {
		user := storage.User{Id: c.Chat().ID, Time: formattedTime}
		if err := user.AddUserToMailing(); err != nil {
			return err
		}
		response := fmt.Sprintf("Рассылка в %s подключена🤙🤣🤣", formattedTime)
		return c.Send(response, createMenu(true))
	}
	return c.Send("Неверный формат времени")
}

func createMenu(isMailing bool) *tele.ReplyMarkup {
	markup := tele.ReplyMarkup{}
	if isMailing {
		markup.Reply(markup.Row(
			markup.Text("Поздравь меня🥳"),
			markup.Text("Ура!🎉"),
			markup.Text("Отключить рассылку📬"),
		))
		return &markup
	} else {
		markup.Reply(markup.Row(
			markup.Text("Поздравь меня🥳"),
			markup.Text("Ура!🎉"),
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
	isMailing, err := storage.GetIfUserIsMailing(c.Chat().ID)
	if err != nil {
		return err
	}

	markup := createMenu(isMailing)
	switch c.Text() {
	case "Ура!🎉":
		return c.Send("Ура!🎉", markup)
	case "Поздравь меня🥳":
		{
			postcardPath, err := storage.GetRandomPostcardPath()
			if err != nil {
				return err
			}
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
			return c.Send("Чтобы подключить рассылку напиши /time *time*, например /time 08:30", markup)
		}
	case "Отключить рассылку📬":
		{
			err = storage.RemoveUserFromMailing(c.Chat().ID)
			if err != nil {
				return err
			}

			return c.Send("Рассылка отключена😈", createMenu(false))
		}
	}

	return nil
}
