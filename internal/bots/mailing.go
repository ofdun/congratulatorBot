package bots

import (
	"CongratulatorBot/internal/download"
	"CongratulatorBot/internal/storage"
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

func EveryMinuteLoop(telegramBot *TelegramBot, errorChan chan error) {
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
		doMailing(telegramBot, errorChan)
	}
}

func doMailing(bot *TelegramBot, errorChan chan error) {
	timeSinceMidnight := getTimeSinceMidnightUTC()
	timeSinceMidnight -= timeSinceMidnight % 60
	telegramUsers, err := storage.GetIDsFromTime(timeSinceMidnight)
	if err != nil {
		errorChan <- err
	}
	if len(telegramUsers) > 0 {
		for _, user := range telegramUsers {
			if err = sendPostcard(bot, user, errorChan); err != nil {
				errorChan <- err
			}
		}
	}
}
