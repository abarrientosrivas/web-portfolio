package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
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
		time.Sleep(1 * time.Second)
		log.Print("HTMX Hit")
		title := r.PostFormValue("Title")
		fmt.Println(title)
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.ExecuteTemplate(w, "test-list-element", Works{Title: title, Description: "manually added"})
		// htmlStr := fmt.Sprintf("<p>%s - manually added</p>", title)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		// tmpl.Execute(w, nil)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/language_selection.html"))
		tmpl.Execute(w, nil)
	}

	h4 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/welcome_es.html"))
		tmpl.Execute(w, nil)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", h1)
	http.HandleFunc("/hit/", h2)
	http.HandleFunc("/initial", h3)
	http.HandleFunc("/lang/es", h4)

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
