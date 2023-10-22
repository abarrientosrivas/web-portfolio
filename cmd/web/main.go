package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type CommonConfig struct {
	RevealHeader bool
	HideHeader   bool
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

type LandingStrings struct {
	WelcomeText string `toml:"WelcomeText"`
	MessageText string `toml:"MessageText"`
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
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/language", LanguageHandler)
	mux.HandleFunc("/legal", LegalHandler)
	mux.HandleFunc("/about", AboutHandler)
	mux.HandleFunc("/work", WorkPage)
	mux.HandleFunc("/contact", ContactHandler)

	log.Print("Server listening on: 127.0.0.1:8000")
	log.Fatal(http.ListenAndServe(":8000", sessionManager.LoadAndSave(mux)))
}

func GetCurrentLanguage(r *http.Request) string {
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

	return language
}

func GetLanguageStrings[T any](language string, configName string) (T, error) {
	var readConfig T
	tomlFilepath := fmt.Sprintf("data/%s/%s.toml", language, configName)

	if _, err := os.Stat(tomlFilepath); os.IsNotExist(err) {
		log.Print("File does not exist:", tomlFilepath)
		log.Print("Trying default.")
		tomlFilepath = fmt.Sprintf("data/en/%s.toml", configName)
		if _, err := os.Stat(tomlFilepath); os.IsNotExist(err) {
			log.Print("File does not exist:", tomlFilepath)
			return readConfig, err
		}
	}

	tomlData, err := os.ReadFile(tomlFilepath)
	if err != nil {
		fmt.Print("Error reading file:", err)
		return readConfig, err
	}

	if _, err := toml.Decode(string(tomlData), &readConfig); err != nil {
		fmt.Print("Error decoding TOML:", err)
		return readConfig, err
	}

	return readConfig, nil
}

func areHeadersAbsent(r *http.Request, headers []string) bool {
	for _, header := range headers {
		if _, ok := r.Header[http.CanonicalHeaderKey(header)]; ok {
			return false
		}
	}
	log.Print("no header")
	return true
}

func IsJustArriving(r *http.Request) bool {
	if areHeadersAbsent(r, []string{"Sec-Fetch-Site", "Origin", "Referer"}) {
		return true
	} else {

		secFetchSite := r.Header.Get("Sec-Fetch-Site")
		if secFetchSite == "same-origin" || secFetchSite == "same-site" {
			return false
		}

		if secFetchSite == "none" || secFetchSite == "cross-site" {
			return true
		}

		origin := r.Header.Get("Origin")
		if origin != "" {
			parsedOrigin, err := url.Parse(origin)
			if err == nil && parsedOrigin.Host != r.Host {
				return true
			} else {
				return false
			}
		}

		referer := r.Header.Get("Referer")
		if referer != "" {
			refererURL, err := url.Parse(referer)
			if err == nil && refererURL.Host != r.Host {
				return true
			} else {
				return false
			}
		}
		return true
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	language := GetCurrentLanguage(r)
	if values, ok := r.Header[http.CanonicalHeaderKey("Referer")]; ok {
		// The header exists.
		log.Printf("%s header exists with value(s): %v\n", "Referer", values)
	} else {
		// The header does not exist.
		log.Printf("%s header does not exist\n", "Referer")
	}

	if values, ok := r.Header[http.CanonicalHeaderKey("Sec-Fetch-Site")]; ok {
		// The header exists.
		log.Printf("%s header exists with value(s): %v\n", "Sec-Fetch-Site", values)
	} else {
		// The header does not exist.
		log.Printf("%s header does not exist\n", "Sec-Fetch-Site")
	}

	commonStrings, err := GetLanguageStrings[CommonStrings](language, "common")
	if err != nil {
		log.Print("Error loading strings", err)
		return
	}

	if IsJustArriving(r) {
		landingStrings, err := GetLanguageStrings[LandingStrings](language, "landing")
		if err != nil {
			log.Print("Error loading strings", err)
			return
		}
		commonConfig := CommonConfig{
			HideHeader: false,
		}
		context := map[string]interface{}{
			"CommonConfig":   commonConfig,
			"CommonStrings":  commonStrings,
			"LandingStrings": landingStrings,
		}

		tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/presentation.html"))
		tmpl.ExecuteTemplate(w, "common", context)
	} else {
		presentationStrings, err := GetLanguageStrings[PresentationStrings](language, "presentation")
		if err != nil {
			log.Print("Error loading strings", err)
			return
		}
		commonConfig := CommonConfig{
			RevealHeader: strings.HasSuffix(r.Header.Get("Referer"), r.URL.String()),
		}

		context := map[string]interface{}{
			"CommonConfig":        commonConfig,
			"CommonStrings":       commonStrings,
			"PresentationStrings": presentationStrings,
		}

		tmpl := template.Must(template.ParseFiles("templates/common.html", "templates/presentation.html"))
		tmpl.ExecuteTemplate(w, "common", context)
	}
}

func LanguageHandler(w http.ResponseWriter, r *http.Request) {
	language := GetCurrentLanguage(r)
	langParameter := r.URL.Query().Get("lang")

	if langParameter != "" {
		language = langParameter
		sessionManager.Put(r.Context(), "language", language)
	}

	commonStrings, err := GetLanguageStrings[CommonStrings](language, "common")
	if err != nil {
		log.Print("Error loading strings", err)
		return
	}
	languageSelectorStrings, err := GetLanguageStrings[LanguageSelectorStrings](language, "language_selector")
	if err != nil {
		log.Print("Error loading strings", err)
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
