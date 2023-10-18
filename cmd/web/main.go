package main

import (
	"log"
	"net/http"
	"text/template"
)

type Works struct {
	Title       string
	Description string
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", WelcomePage)
	http.HandleFunc("/home", HomePage)
	http.HandleFunc("/language", LanguagePage)
	http.HandleFunc("/legal", LegalPage)
	http.HandleFunc("/about", AboutPage)
	http.HandleFunc("/work", WorkPage)
	http.HandleFunc("/contact", ContactPage)

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		print(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		print(err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func LanguagePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/plain_htmls/language.html"))
	tmpl.Execute(w, nil)
}

func LegalPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/plain_htmls/language.html"))
	tmpl.Execute(w, nil)
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/plain_htmls/language.html"))
	tmpl.Execute(w, nil)
}

func WorkPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/plain_htmls/language.html"))
	tmpl.Execute(w, nil)
}

func ContactPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/plain_htmls/language.html"))
	tmpl.Execute(w, nil)
}
