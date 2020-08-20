package main

import (
	"log"
	"net/http"
)

var LdapConfig *Config

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "register.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	secret := req.FormValue("secret")

	if LdapConfig.Secret != "" && LdapConfig.Secret != secret {
		res.Write([]byte("Get a load of this guy, not knowing the secret code"))
		return
	}
	//insert into LDAP
	log.Printf("Got %v %v %v %v\n", username, password, email, secret)
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
	LdapConfig, _ = LoadConfig()
	log.Println("Loaded config")
	http.HandleFunc("/register", signupPage)
	http.HandleFunc("/", homePage)
	log.Println("Guildgate starting")
	http.ListenAndServe(":8080", nil)

}
