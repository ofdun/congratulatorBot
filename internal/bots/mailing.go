package bots

import (
	"CongratulatorBot/internal/download"
	"CongratulatorBot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"math"
	"os"
	"strconv"
	"time"
)

func getTimeSinceMidnightUTC() int {
	return int(time.Now().Unix()) % 86400
}

func isClose(f int, s int, approx int) bool {
	if math.Abs(float64(f-s)) < float64(approx) {
		return true
	}
	return false
}

func EveryMinuteLoop(discordBot *discordgo.Session, telegramBot *TelegramBot, errorChan chan error) {
	files, err := os.ReadDir("./internal/storage/postcards/")
	if err != nil {
		errorChan <- err
	}
	if len(files) == 0 {
		download.NightlyPostcardDownload(errorChan)
	}

	for {
		time.Sleep(time.Minute)

		offset, err := strconv.Atoi(os.Getenv("TIME_OFFSET"))
		if err != nil {
			errorChan <- err
		}
		timeInUTC3 := (getTimeSinceMidnightUTC() + offset*3600) % 86400
		if isClose(timeInUTC3, 150, 59) {
			download.NightlyPostcardDownload(errorChan)
		}
		doMailing(discordBot, telegramBot, errorChan)
	}
}

func doMailing(discordBot *discordgo.Session, telegramBot *TelegramBot, errorChan chan error) {
	timeSinceMidnight := getTimeSinceMidnightUTC()
	timeSinceMidnight -= timeSinceMidnight % 60
	telegramUsers, err := storage.GetIDsFromTime(timeSinceMidnight, false)
	if err != nil {
		errorChan <- err
	}
	if len(telegramUsers) > 0 {
		for _, user := range telegramUsers {
			if err = sendPostcard(telegramBot, user, errorChan); err != nil {
				errorChan <- err
			}
		}
	}

	discordChannels, err := storage.GetIDsFromTime(timeSinceMidnight, true)
	if len(discordChannels) > 0 {
		for _, channelID := range discordChannels {
			sendDiscordPostcard(discordBot, strconv.Itoa(int(channelID)))
		}
	}
}
