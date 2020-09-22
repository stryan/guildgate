package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/securecookie"
)

var Conf *Config
var tpl *template.Template
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Println("GET /register")
		u := getUserName(req)
		if u != "" {
			http.Redirect(res, req, "/", 302)
		} else {
			data := struct {
				Title    string
				Username string
				LoggedIn bool
			}{
				"Register",
				"Unregistered",
				false,
			}
			tpl.ExecuteTemplate(res, "register", data)
		}
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	secret := req.FormValue("secret")

	if Conf.Secret != "" && Conf.Secret != secret {
		//Checking it as a token
		err := validateToken(secret)
		if err != nil {
			log.Printf("Bad secret entered: %v\n", err)
			res.Write([]byte("Get a load of this guy, not knowing the secret code"))
			return
		}
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

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Println("GET /login")
		u := getUserName(req)
		if u != "" {
			http.Redirect(res, req, "/", 302)
		} else {
			data := struct {
				Title    string
				Username string
				LoggedIn bool
			}{
				"Login",
				"Unregistered",
				false,
			}
			tpl.ExecuteTemplate(res, "login", data)
		}
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	log.Printf("Attempting login for user %v\n", username)
	err := loginLDAPAccount(username, password)
	if err != nil {
		log.Printf("Error logging in user %v: %v\n", username, err)
		res.Write([]byte("Error logging in. Incorrect password?"))
		return
	} else {
		setSession(username, res)
		http.Redirect(res, req, "/", 302)
		return
	}
}

func logoutPage(res http.ResponseWriter, req *http.Request) {
	clearSession(res)
	http.Redirect(res, req, "/", 302)
}

func tokenPage(res http.ResponseWriter, req *http.Request) {
	u := getUserName(req)
	if u == "" {
		http.Redirect(res, req, "/", 302)
	}
	token, err := generateToken(u)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		tpl.ExecuteTemplate(res, "error", nil)
	}
	data := struct {
		Title    string
		Username string
		LoggedIn bool
		Token    string
	}{
		"Token Generation",
		u,
		true,
		token,
	}
	tpl.ExecuteTemplate(res, "token", data)
}

func homePage(res http.ResponseWriter, req *http.Request) {
	u := getUserName(req)
	active := false
	uname := "Unregistered"
	if u != "" {
		uname = u
		active = true
	}
	data := struct {
		Title    string
		Username string
		LoggedIn bool
	}{
		"Index",
		uname,
		active,
	}

	tpl.ExecuteTemplate(res, "index", data)
}

func main() {
	Conf, _ = LoadConfig()
	log.Println("Loaded config")
	http.HandleFunc("/register", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/token", tokenPage)
	http.HandleFunc("/", homePage)
	log.Printf("Registering templates from %v/\n", Conf.TplPath)
	tpl = template.Must(template.ParseGlob(Conf.TplPath + "/*.html"))
	log.Printf("Guildgate starting on %v\n", Conf.Port)
	var err error
	if Conf.Tls {
		log.Printf("Starting TLS\n")
		if Conf.Cert == "" {
			log.Fatalf("Need to specify a certificate if using TLS!\n")
		} else if Conf.Key == "" {
			log.Fatalf("Need to specify a private key is usingTLS!\n")
		} else {
			err = http.ListenAndServeTLS(":"+Conf.Port, Conf.Cert, Conf.Key, nil)
		}
	} else {
		log.Printf("Starting unencrypted\n")
		err = http.ListenAndServe(":"+Conf.Port, nil)
	}
	if err != nil {
		log.Printf("HTTP server failed with %v\n", err)
	}
}
