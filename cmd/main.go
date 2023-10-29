package main

import (
	"telegramBot/internal/download"
	"telegramBot/internal/parser"
)

func main() {
	newParser := parser.NewParser(
		"https://3d-galleru.ru/archive/cat/kalendar-42/")

	err := newParser.GetHTML()
	if err != nil {
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
		err = newParser.GetPostcardHref(&postcards[i])
		if err != nil {
			panic(err)
		}
	}

	for i := range postcards {
		err = download.PostcardDownload("internal/storage/postcards/", &postcards[i])
		if err != nil {
			panic(err)
		}
	}
}
