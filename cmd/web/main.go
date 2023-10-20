package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type CommonConfig struct {
	ShowHeader bool
	HideHeader bool
}

type CommonStrings struct {
	HomeText          string `toml:"HomeText"`
	AboutText         string `toml:"AboutText"`
	WorkText          string `toml:"WorkText"`
	ContactText       string `toml:"ContactText"`
	LanguageText      string `toml:"LanguageText"`
	CopyrightText     string `toml:"CopyrightText"`
	LicenseText       string `toml:"LicenseText"`
	PrivacyPolicyText string `toml:"PrivacyPolicyText"`
}

type PresentationStrings struct {
	WelcomeText      string `toml:"WelcomeText"`
	PresentationText string `toml:"PresentationText"`
	MessageText      string `toml:"MessageText"`
}

type LanguageSelectorStrings struct {
	LanguageSelectionText string `toml:"LanguageSelectionText"`
}

var sessionManager *scs.SessionManager

func main() {
	sessionManager = scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", LandingHandler)
	mux.HandleFunc("/welcome", WelcomeHandler)
	mux.HandleFunc("/language", LanguageHandler)
	mux.HandleFunc("/legal", LegalHandler)
	mux.HandleFunc("/about", AboutHandler)
	mux.HandleFunc("/work", WorkPage)
	mux.HandleFunc("/contact", ContactHandler)

	log.Print("Server listening on: 127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", sessionManager.LoadAndSave(mux)))
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	var language string
	if sessionManager.GetString(r.Context(), "language") == "" {
		preferredLanguage := r.Header.Get("Accept-Language")
		if preferredLanguage == "" {
			preferredLanguage = "en"
		}
		preferredLanguage = preferredLanguage[:2]
		sessionManager.Put(r.Context(), "language", preferredLanguage)
		language = preferredLanguage
	} else {
		language = sessionManager.GetString(r.Context(), "language")
	}

	commonFilePath := fmt.Sprintf("data/%s/common.toml", language)
	presentationFilePath := fmt.Sprintf("data/%s/presentation.toml", language)

	if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", commonFilePath)
		log.Println("Trying default.")
		commonFilePath = "data/en/common.toml"
		if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", commonFilePath)
			return
		}
	}

	if _, err := os.Stat(presentationFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", presentationFilePath)
		log.Println("Trying default.")
		presentationFilePath = "data/en/presentation.toml"
		if _, err := os.Stat(presentationFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", presentationFilePath)
			return
		}
	}

	tomlData, err := os.ReadFile(commonFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var commonStrings CommonStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &commonStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	tomlData, err = os.ReadFile(presentationFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var presentationStrings PresentationStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &presentationStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	commonConfig := CommonConfig{
		HideHeader: true,
	}

	context := map[string]interface{}{
		"CommonConfig":        commonConfig,
		"CommonStrings":       commonStrings,
		"PresentationStrings": presentationStrings,
	}

	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/landing.html"))
	tmpl.ExecuteTemplate(w, "common", context)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	preferredLanguage := r.Header.Get("Accept-Language")
	if preferredLanguage == "" {
		preferredLanguage = "en"
	}
	preferredLanguage = preferredLanguage[:2]

	commonFilePath := fmt.Sprintf("data/%s/common.toml", preferredLanguage)
	presentationFilePath := fmt.Sprintf("data/%s/presentation.toml", preferredLanguage)

	if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", commonFilePath)
		log.Println("Trying default.")
		commonFilePath = "data/en/common.toml"
		if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", commonFilePath)
			return
		}
	}

	if _, err := os.Stat(presentationFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", presentationFilePath)
		log.Println("Trying default.")
		presentationFilePath = "data/en/presentation.toml"
		if _, err := os.Stat(presentationFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", presentationFilePath)
			return
		}
	}

	tomlData, err := os.ReadFile(commonFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var commonStrings CommonStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &commonStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	tomlData, err = os.ReadFile(presentationFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var presentationStrings PresentationStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &presentationStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	commonConfig := CommonConfig{
		ShowHeader: true,
	}

	context := map[string]interface{}{
		"CommonConfig":        commonConfig,
		"CommonStrings":       commonStrings,
		"PresentationStrings": presentationStrings,
	}

	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/presentation.html"))
	tmpl.ExecuteTemplate(w, "common", context)
}

func LanguageHandler(w http.ResponseWriter, r *http.Request) {
	var language string
	if sessionManager.GetString(r.Context(), "language") == "" {
		preferredLanguage := r.Header.Get("Accept-Language")
		if preferredLanguage == "" {
			preferredLanguage = "en"
		}
		preferredLanguage = preferredLanguage[:2]
		sessionManager.Put(r.Context(), "language", preferredLanguage)
		language = preferredLanguage
	} else {
		language = sessionManager.GetString(r.Context(), "language")
	}

	langParam := r.URL.Query().Get("lang")

	if langParam != "" {
		language = langParam
		sessionManager.Put(r.Context(), "language", language)
	}

	commonFilePath := fmt.Sprintf("data/%s/common.toml", language)
	languageSelectorFilePath := fmt.Sprintf("data/%s/language_selector.toml", language)

	if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", commonFilePath)
		log.Println("Trying default.")
		commonFilePath = "data/en/common.toml"
		if _, err := os.Stat(commonFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", commonFilePath)
			return
		}
	}

	if _, err := os.Stat(languageSelectorFilePath); os.IsNotExist(err) {
		log.Println("File does not exist:", languageSelectorFilePath)
		log.Println("Trying default.")
		languageSelectorFilePath = "data/en/language_selector.toml"
		if _, err := os.Stat(languageSelectorFilePath); os.IsNotExist(err) {
			log.Println("File does not exist:", languageSelectorFilePath)
			return
		}
	}

	tomlData, err := os.ReadFile(commonFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var commonStrings CommonStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &commonStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	tomlData, err = os.ReadFile(languageSelectorFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var languageSelectorStrings LanguageSelectorStrings

	// Parse the TOML data into the Config struct.
	if _, err := toml.Decode(string(tomlData), &languageSelectorStrings); err != nil {
		fmt.Println("Error decoding TOML:", err)
		return
	}

	context := map[string]interface{}{
		"CommonConfig":            CommonConfig{},
		"CommonStrings":           commonStrings,
		"LanguageSelectorStrings": languageSelectorStrings,
	}

	tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/language_selector.html"))
	tmpl.ExecuteTemplate(w, "common", context)
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
