package test

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

var artists []Artist2
var locindex LocationIndex
var dateindex DateIndex
var relindex RelationIndex
var Wrong string

func RunAllApi() {
	// Получение данных из API
	api1 := "https://groupietrackers.herokuapp.com/api/locations"
	api2 := "https://groupietrackers.herokuapp.com/api/artists"
	api3 := "https://groupietrackers.herokuapp.com/api/dates"
	api4 := "https://groupietrackers.herokuapp.com/api/relation"

	// Получение данных об индексе местоположений
	respLocation, err := http.Get(api1)
	if err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}
	defer respLocation.Body.Close()

	if err := json.NewDecoder(respLocation.Body).Decode(&locindex); err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}

	// Получение данных об артистах
	respArtists, err := http.Get(api2)
	if err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}
	defer respArtists.Body.Close()

	if err := json.NewDecoder(respArtists.Body).Decode(&artists); err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}

	respDate, err := http.Get(api3)
	if err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}
	defer respDate.Body.Close()

	if err := json.NewDecoder(respDate.Body).Decode(&dateindex); err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}

	respRelation, err := http.Get(api4)
	if err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}
	defer respRelation.Body.Close()

	if err := json.NewDecoder(respRelation.Body).Decode(&relindex); err != nil {
		Wrong = "Error"
		log.Println("ERROR")
		return
	}

	// Создание структуры Artist1, объединяющей данные из artists и index
	for i, artist := range artists {
		for _, locData := range locindex.Index {
			if locData.ID == artist.ID {
				artists[i].Location = locData
				break
			}
		}
	}

	for i, artist := range artists {
		for _, locData := range dateindex.Index {
			if locData.ID == artist.ID {
				artists[i].Date = locData
				break
			}
		}
	}

	for i, artist := range artists {
		for _, locData := range relindex.Index {
			if locData.ID == artist.ID {
				artists[i].Relation = locData
				break
			}
		}
	}
}

func ArtistFunc(w http.ResponseWriter, r *http.Request) {
	if Wrong != "" {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	if r.URL.Path != "/" {
		tmpl, _ := template.ParseFiles("client/404.html")
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, nil)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("client/artists.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
}

func ArtistByIDFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := path.Base(r.URL.Path)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/404.html")
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, nil)
		return
	}

	var artist Artist2
	for _, a := range artists {
		if a.ID == idInt {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		tmpl, _ := template.ParseFiles("client/404.html")
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, nil)
		return
	}

	tmpl, err := template.ParseFiles("client/artist.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, artist)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
}

func SearchArtists(w http.ResponseWriter, r *http.Request) {

	flag := true

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var ArtistsRes []Artist2

	query := strings.ToLower(r.URL.Query().Get("query"))
	if len(query) > 1 {
		for _, val := range artists {
			if strings.Contains(strings.ToLower(val.Name), query) ||
				strings.Contains(strings.ToLower(val.FirstAlbum), query) ||
				strings.Contains(strings.ToLower(strconv.Itoa(val.CreationDate)), query) {
				ArtistsRes = append(ArtistsRes, val)
				flag = false
			} else {
				flag = true
			}

			for _, member := range val.Members {
				if strings.Contains(strings.ToLower(member), query) && flag {
					ArtistsRes = append(ArtistsRes, val)
					break
				}
			}

			for _, location := range val.Location.Locations {
				if strings.Contains(strings.ToLower(location), query) {
					ArtistsRes = append(ArtistsRes, val)
					break
				}

			}
		}

	} else {
		for _, val := range artists {
			if strings.Contains(strings.ToLower(val.Name), query) {
				ArtistsRes = append(ArtistsRes, val)
			}
		}
	}

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

		w.Header().Set("Content-Type", "application/json")
		encodedData, err := json.Marshal(ArtistsRes)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных в JSON", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(encodedData)
		if err != nil {
			http.Error(w, "Ошибка при записи данных в ответ", http.StatusInternalServerError)
			return
		}
	} else {
		if ArtistsRes == nil {
			tmpl, err := template.ParseFiles("client/404.html")
			if err != nil {
				return
			}
			err = tmpl.Execute(w, nil)

		} else {
			tmpl, err := template.ParseFiles("client/artists.html")
			if err != nil {
				return
			}
			err = tmpl.Execute(w, ArtistsRes)
		}
	}

}
