package main

import (
	"log"
	"net/http"
)

func setSession(uname string, res http.ResponseWriter) {
	value := map[string]string{
		"name": uname,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}

func getUserName(req *http.Request) (uname string) {
	if cookie, err := req.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			uname = cookieValue["name"]
		}
	}
	return uname
}

func clearSession(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}

func signup(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	secret := req.FormValue("secret")

	if Conf.Secret != "" && Conf.Secret != secret {
		//Checking it as a token
		_, err := validateToken(secret)
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

func login(res http.ResponseWriter, req *http.Request) {
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

func logout(res http.ResponseWriter, req *http.Request) {
	clearSession(res)
}
