package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
)

type PresentationConfig struct {
	WelcomeText      string `toml:"WelcomeText"`
	PresentationText string `toml:"PresentationText"`
	MessageText      string `toml:"MessageText"`
}

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
	filePath := "data/en/presentation.toml"

	// Check if the file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", filePath)
		return
	}

	tomlData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var config PresentationConfig

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &config); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	data := struct {
		ShowHeader       bool
		HideHeader       bool
		WelcomeText      string
		PresentationText string
		MessageText      string
	}{
		ShowHeader:       true,
		WelcomeText:      config.WelcomeText,
		PresentationText: config.PresentationText,
		MessageText:      config.MessageText,
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
