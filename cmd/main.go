package main

import (
	"fmt"
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

	for _, h := range holidays {
		fmt.Println(h)
	}

	//html := newParser.HTML
	//
	//println(html)
}
