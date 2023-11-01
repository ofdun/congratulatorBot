package main

import (
	_ "github.com/lib/pq"
)

func main() {
	//postcardsDatabase := storage.NewDatabase()
	//postcardsStorage := storage.NewPostcardsPostgresStorage(postcardsDatabase)
	//
	//newParser := parser.NewParser(
	//	"https://3d-galleru.ru/archive/cat/kalendar-42/")
	//
	//if err := newParser.GetHTML(); err != nil {
	//	panic(err)
	//}
	//
	//holidays, err := newParser.GetHolidays()
	//if err != nil {
	//	panic(err)
	//}
	//
	//postcards, err := newParser.GetPostcardsPages(holidays[0])
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Updating postcards ( adding links )
	//for i := range postcards {
	//	if err = newParser.GetPostcardHref(&postcards[i]); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//for i := range postcards {
	//	if err = download.PostcardDownload("internal/storage/postcards/", &postcards[i]); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//for i := range postcards {
	//	if err = postcardsStorage.AddPostcard(&postcards[i]); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//postcardsFromStorage, err := postcardsStorage.GetPostcardsFromStorage()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(postcardsFromStorage)
}
