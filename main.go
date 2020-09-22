package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var Conf *Config
var tpl *template.Template
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func main() {
	Conf, _ = LoadConfig()
	log.Println("Loaded config")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/register", signupPage).Methods("GET")
	router.HandleFunc("/register", signup).Methods("POST")
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logoutPage).Methods("GET")
	router.HandleFunc("/token", tokenPage).Methods("GET")
	router.HandleFunc("/passwordreset", resetPageFront).Methods("GET")
	router.HandleFunc("/passwordreset", resetLookup).Methods("POST")
	router.HandleFunc("/passwordresetform", resetPageBack).Methods("GET")
	router.HandleFunc("/passwordresetform", reset).Methods("POST")
	router.HandleFunc("/resetsuccess", resetSuccessPage).Methods("GET")
	router.HandleFunc("/reseterror", resetErrorPage).Methods("GET")
	log.Printf("Registering templates from %v/\n", Conf.TplPath)
	tpl = template.Must(template.ParseGlob(Conf.TplPath + "/*"))
	log.Printf("Guildgate starting on %v\n", Conf.Port)
	var err error
	if Conf.Tls {
		log.Printf("Starting TLS\n")
		if Conf.Cert == "" {
			log.Fatalf("Need to specify a certificate if using TLS!\n")
		} else if Conf.Key == "" {
			log.Fatalf("Need to specify a private key is usingTLS!\n")
		} else {
			err = http.ListenAndServeTLS(":"+Conf.Port, Conf.Cert, Conf.Key, router)
		}
	} else {
		log.Printf("Starting unencrypted\n")
		err = http.ListenAndServe(":"+Conf.Port, router)
	}
	if err != nil {
		log.Printf("HTTP server failed with %v\n", err)
	}
}
