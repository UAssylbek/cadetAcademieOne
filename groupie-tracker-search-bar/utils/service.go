package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

type Location struct {
	City []string `json:"city"`
}

type DatesLocationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type LocationData struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type CombinedData struct {
	Artists   Artist
	Relation  DatesLocationsData
	Dates     DatesData
	Locations LocationData
}

var Artists []Artist
var Relation []DatesLocationsData
var Dates []DatesData
var Locations []LocationData

func ArtistFunc(w http.ResponseWriter, r *http.Request) {

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

	api := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(api)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&Artists)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	tmpl, err := template.ParseFiles("client/artists.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, Artists)
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

	api := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(api)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&Artists)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	var artist Artist
	for _, a := range Artists {
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

	url1 := "https://groupietrackers.herokuapp.com/api/relation/" + id
	resp1, err := http.Get(url1)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	defer resp1.Body.Close()

	var relationn DatesLocationsData
	jsonData, err := ioutil.ReadAll(resp1.Body)
	err = json.Unmarshal(jsonData, &relationn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	url2 := "https://groupietrackers.herokuapp.com/api/dates/" + id
	resp2, err := http.Get(url2)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	defer resp2.Body.Close()

	var datess DatesData
	jsonData, err = ioutil.ReadAll(resp2.Body)
	err = json.Unmarshal(jsonData, &datess)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	url3 := "https://groupietrackers.herokuapp.com/api/locations/" + id
	resp3, err := http.Get(url3)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	defer resp3.Body.Close()

	var locationss LocationData
	jsonData, err = ioutil.ReadAll(resp3.Body)
	err = json.Unmarshal(jsonData, &locationss)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}


	var locationsss Location

	jsonData2, err := ioutil.ReadAll(resp3.Body)
	err = json.Unmarshal(jsonData2, &locationsss)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	jsonData1, err := json.Marshal(locationsss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON клиенту
	w.Write(jsonData1)


	tmpl, err := template.ParseFiles("client/artist.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
	combinedData := CombinedData{
		Artists:   artist,
		Relation:  relationn,
		Dates:     datess,
		Locations: locationss,
	}

	err = tmpl.Execute(w, combinedData)
	if err != nil {
		tmpl, _ := template.ParseFiles("client/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
}
