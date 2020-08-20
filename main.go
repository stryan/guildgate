package main

import (
	"log"
	"net/http"
)

var Conf *Config

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "register.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	secret := req.FormValue("secret")

	if Conf.Secret != "" && Conf.Secret != secret {
		log.Printf("Bad secret entered\n")
		res.Write([]byte("Get a load of this guy, not knowing the secret code"))
		return
	}
	//insert into LDAP
	log.Printf("Attempting to create account for %v", username)
	err := createLDAPAccount(username, password, email)
	if err == nil {
		res.Write([]byte("User created!"))
		return
	} else {
		res.Write([]byte("Failure to create account"))
		return
	}
}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func main() {
	Conf, _ = LoadConfig()
	log.Println("Loaded config")
	http.HandleFunc("/register", signupPage)
	http.HandleFunc("/", homePage)
	log.Printf("Guildgate starting on %v\n", Conf.Port)
	http.ListenAndServe(":"+Conf.Port, nil)

}
