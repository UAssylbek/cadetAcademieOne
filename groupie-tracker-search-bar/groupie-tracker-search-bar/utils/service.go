package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

type LocationIndex struct {
	Index []LocationData `json:"index"`
}
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

func SearchArtists(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var ArtistsRes []Artist

	query := strings.ToLower(r.URL.Query().Get("query"))
	index, err := AllLocations()
	artists, err := AllArtists()

	if err != nil {
		http.Error(w, "ошибка при получений данных", http.StatusInternalServerError)
	}

	var locat []int

	for _, i := range index.Index {
		for _, j := range i.Locations {
			if strings.Contains(strings.ToLower(j), query) {
				locat = append(locat, i.ID)
			}
		}
	}

	for _, val := range artists {
		if strings.Contains(strings.ToLower(val.Name), query) || 
		strings.Contains(strings.ToLower(val.FirstAlbum), query) || 
		strings.Contains(strings.ToLower(strconv.Itoa(val.CreationDate)), query) || 
		IfDublicate(locat, val.ID) {
			ArtistsRes = append(ArtistsRes, val)
		}
	}

	for _, member := range artists {
		for _, val := range member.Members {
			if strings.Contains(strings.ToLower(val), query) {
				ArtistsRes = append(ArtistsRes, member)
				break
			}
		}
	}

	if ArtistsRes ==  nil {
		tmpl, err := template.ParseFiles("client/artistsnil.html")
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

func AllLocations() (LocationIndex, error) {
	api := "https://groupietrackers.herokuapp.com/api/locations"

	resp, err := http.Get(api)
	if err != nil {
		return LocationIndex{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationIndex{}, err
	}
	var index LocationIndex
	err = json.Unmarshal(data, &index)
	if err != nil {
		return LocationIndex{}, err
	}

	return index, nil
}

func AllArtists() ([]Artist, error) {
	api := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(api)
	if err != nil {
		return []Artist{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Artist{}, err
	}
	var artists []Artist
	err = json.Unmarshal(data, &artists)
	if err != nil {
		return []Artist{}, err
	}

	return artists, nil
}

func IfDublicate(slice []int, id int) bool {
	for _, i := range slice {
		if i == id {
			return true
		}
	}
	return false
}