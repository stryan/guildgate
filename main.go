package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var Conf *Config
var tpl *template.Template
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
var passwordTokenSet map[string]bool
var GitCommit string

func main() {
	Conf, _ = LoadConfig()
	version := GitCommit
	log.Println("Loaded config")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/register", signupPage).Methods("GET")
	router.HandleFunc("/register", signup).Methods("POST")
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logoutPage).Methods("GET")
	router.HandleFunc("/token", tokenPage).Methods("GET")
	router.HandleFunc("/profile/view", profilePage).Methods("GET")
	router.HandleFunc("/profile/edit", profileEditPage).Methods("GET")
	router.HandleFunc("/profile/edit", profileEdit).Methods("POST")
	router.HandleFunc("/minecraft", minecraftPage).Methods("GET")
	router.HandleFunc("/minecraft/link", minecraftLink).Methods("POST")
	router.HandleFunc("/minecraft/link/success", minecraftLinkSuccessPage).Methods("GET")
	router.HandleFunc("/minecraft/link/error", minecraftLinkErrorPage).Methods("GET")
	router.HandleFunc("/reset", resetPageFront).Methods("GET")
	router.HandleFunc("/reset", resetLookup).Methods("POST")
	router.HandleFunc("/reset/form", resetPageBack).Methods("GET")
	router.HandleFunc("/reset/form", reset).Methods("POST")
	router.HandleFunc("/reset/success", resetSuccessPage).Methods("GET")
	router.HandleFunc("/reset/error", resetErrorPage).Methods("GET")
	router.HandleFunc("/change", changePageFront).Methods("GET")
	router.HandleFunc("/change", change).Methods("POST")
	router.HandleFunc("/change/success", changeSuccessPage).Methods("GET")
	router.HandleFunc("/change/error", changeSuccessPage).Methods("GET")

	log.Printf("Registering templates from %v/\n", Conf.TplPath)
	tpl = template.Must(template.ParseGlob(Conf.TplPath + "/*"))
	if Conf.UserTplPath != "" {
		log.Printf("Registering user templates from %v/\n", Conf.UserTplPath)
		tpl = template.Must(tpl.ParseGlob(Conf.UserTplPath + "/*"))
	}
	log.Println("Performing LDAP checks")
	log.Println("Loading max employeeNumber for account creation")
	i, err := findLDAPMaxID()
	if err != nil {
		log.Printf("WARN: Unable to calculate max employeeNumber: %v\n", err)
	} else {
		Conf.MaxID = i
		log.Printf("Max employeeNumber set to %v\n", Conf.MaxID)
	}
	passwordTokenSet = make(map[string]bool)
	log.Printf("Guildgate v%v starting on %v\n", version, Conf.Port)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)
	go func() {
		for {
			<-sigs
			log.Println("reloading templates on SIGUSR1")
			tpl = template.Must(template.ParseGlob(Conf.TplPath + "/*"))
			if Conf.UserTplPath != "" {
				tpl = template.Must(tpl.ParseGlob(Conf.UserTplPath + "/*"))
			}
		}
	}()
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
