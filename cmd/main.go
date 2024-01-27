package main

import (
	"CongratulatorBot/internal/bots"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		panic(err)
	}

	errorChan := make(chan error, 1)

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramBot, err := bots.NewTelegramBot(telegramToken)
	if err != nil {
		panic(err)
	}

	go bots.StartTelegramBot(telegramBot, errorChan)
	go bots.EveryMinuteLoop(telegramBot, errorChan)

	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	go bots.StartDiscordBot(discordToken, errorChan)

	select {
	case err = <-errorChan:
		panic(err)
	}
}
