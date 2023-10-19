package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", LandingHandler)
	http.HandleFunc("/welcome", WelcomeHandler)
	http.HandleFunc("/language", LanguageHandler)
	http.HandleFunc("/legal", LegalHandler)
	http.HandleFunc("/about", AboutHandler)
	http.HandleFunc("/work", WorkPage)
	http.HandleFunc("/contact", ContactHandler)

	log.Print("Server listening on: 127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ShowHeader bool
		HideHeader bool
	}{
		HideHeader: true,
	}
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/landing.html"))
	tmpl.ExecuteTemplate(w, "common", data)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ShowHeader bool
		HideHeader bool
	}{
		ShowHeader: true,
	}
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/presentation.html"))
	tmpl.ExecuteTemplate(w, "common", data)
}

func LanguageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/language_selector.html"))
	tmpl.ExecuteTemplate(w, "common", nil)
}

func LegalHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/legal_page.html"))
	tmpl.ExecuteTemplate(w, "common", nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/resume.html"))
	tmpl.ExecuteTemplate(w, "common", nil)
}

func WorkPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/works_index.html"))
	tmpl.ExecuteTemplate(w, "common", nil)
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/contact_info.html"))
	tmpl.ExecuteTemplate(w, "common", nil)
}
