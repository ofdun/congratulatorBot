package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"telegramBot/internal/bots"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		panic(err)
	}

	//telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	//telegramBot, err := bots.NewTelegramBot(telegramToken)
	//if err != nil {
	//	panic(err)
	//}
	//
	//go bots.StartTelegramBot(telegramBot)
	//go bots.EveryMinuteLoop(telegramBot)
	go bots.StartDiscordBot()

	select {}
}
