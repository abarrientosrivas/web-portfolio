package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Works struct {
	Title       string
	Description string
}

func main() {
	fmt.Println("hello world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		works := map[string][]Works{
			"Works": {
				{Title: "Proj A", Description: "The proy A is cool"},
				{Title: "Proj B", Description: "B is for bad...."},
			},
		}
		tmpl.Execute(w, works)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		log.Print("HTMX Hit")
		log.Print(r.Header.Get("HX-Request"))
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", h1)
	http.HandleFunc("/hit/", h2)

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
