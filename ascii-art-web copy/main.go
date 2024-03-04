package main

import (
	"ascii/ascii"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var OutputText string

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ascii-art", peopleHandler)
	http.HandleFunc("/health", healtCheckHandler)

	log.Println("serve start listening on port 8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "404.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		getPeople(w, r)
	case http.MethodPost:
		postPerson(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", OutputText)
}

func postPerson(w http.ResponseWriter, r *http.Request) {
	Text := r.FormValue("fname")
	Style := r.FormValue("styleName")
	wrong := ""
	OutputText, wrong = ascii.AsciiImpl(Text, Style)
	if wrong == "no Latin" {
		tmpl, _ := template.ParseFiles("static/400.html")
		w.WriteHeader(http.StatusBadRequest)
		tmpl.Execute(w, nil)
		return
	} else if wrong == "no banner" {
		tmpl, _ := template.ParseFiles("static/500.html")
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func healtCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "static/404.html")
		return
	}

	fmt.Fprint(w, "http web-server works correctly")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "static/404.html")
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}
