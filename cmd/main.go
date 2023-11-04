package main

import (
	_ "github.com/lib/pq"
	"os"
	"telegramBot/internal/bots"
	"telegramBot/internal/download"
	"telegramBot/internal/parser"
	"telegramBot/internal/storage"
)

func main() {
	postcardsDatabase := storage.NewDatabase()
	postcardsStorage := storage.NewPostcardsPostgresStorage(postcardsDatabase)
	if err := postcardsStorage.ClearDatabase(); err != nil {
		panic(err)
	}

	newParser := parser.NewParser(
		"https://3d-galleru.ru/archive/cat/kalendar-42/")

	if err := newParser.GetHTML(); err != nil {
		panic(err)
	}

	holidays, err := newParser.GetHolidays()
	if err != nil {
		panic(err)
	}

	postcards, err := newParser.GetPostcardsPages(holidays[0])

	if err != nil {
		panic(err)
	}

	// Updating postcards ( adding links )
	for i := range postcards {
		if err = newParser.GetPostcardHref(&postcards[i]); err != nil {
			panic(err)
		}
	}

	if err = os.RemoveAll("internal/storage/postcards/"); err != nil {
		panic(err)
	}
	if err = os.Mkdir("internal/storage/postcards/", 777); err != nil {
		panic(err)
	}

	for i := range postcards {
		if err = download.PostcardDownload("internal/storage/postcards/", &postcards[i]); err != nil {
			panic(err)
		}
	}

	for i := range postcards {
		if err = postcardsStorage.AddPostcardToStorage(&postcards[i]); err != nil {
			panic(err)
		}
	}

	token := "TOKEN"
	bot, err := bots.NewTelegramBot(token)
	if err != nil {
		panic(err)
	}

	updates, err := bots.SetupUpdates(&bot)
	if err != nil {
		panic(err)
	}

	go bots.StartBotPolling(&bot, updates)

	select {}
}
