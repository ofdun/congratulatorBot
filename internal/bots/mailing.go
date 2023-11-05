package bots

import (
	"math"
	"os"
	"strconv"
	"telegramBot/internal/download"
	"telegramBot/internal/storage"
	"time"
)

func getTimeSinceMidnightUTC() int {
	return int(time.Now().Unix()) % 86400
}

func isClose(f int, s int, aprox int) bool {
	if math.Abs(float64(f-s)) < float64(aprox) {
		return true
	}
	return false
}

func EveryMinuteLoop(telegramBot *TelegramBot) {
	for {
		time.Sleep(time.Minute)

		offset, err := strconv.Atoi(os.Getenv("TIME_OFFSET"))
		if err != nil {
			panic(err)
		}
		timeInUTC3 := (getTimeSinceMidnightUTC() + offset*3600) % 86400
		if isClose(timeInUTC3, 150, 59) {
			if err := download.NightlyPostcardDownload(); err != nil {
				panic(err)
			}
		}
		doMailing(telegramBot)
	}
}

func doMailing(bot *TelegramBot) {
	timeSinceMidnight := getTimeSinceMidnightUTC()
	timeSinceMidnight -= timeSinceMidnight % 60
	telegramUsers, err := storage.GetIDsFromTime(timeSinceMidnight)
	if err != nil {
		panic(err)
	}
	if len(telegramUsers) > 0 {
		for _, user := range telegramUsers {
			if err = sendPostcard(bot, user); err != nil {
				panic(err)
			}
		}
	}
}
