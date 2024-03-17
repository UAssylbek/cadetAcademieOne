package main

import (
	"log"
	"net/http"
	"test/utils"
)

func main() {

	http.HandleFunc("/", test.ArtistFunc)
	http.HandleFunc("/artists/", test.ArtistByIDFunc) // Добавляем новый обработчик для запросов по ID артиста
	http.HandleFunc("/search", test.SearchArtists)
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("client"))))
	log.Println("serve start listening on port 8080")
	errr := http.ListenAndServe("0.0.0.0:8080", nil)
	if errr != nil {
		log.Fatal(errr)
	}
}
