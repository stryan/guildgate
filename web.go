package main

import (
	"log"
	"net/http"
)

func signupPage(res http.ResponseWriter, req *http.Request) {
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

func loginPage(res http.ResponseWriter, req *http.Request) {
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

func logoutPage(res http.ResponseWriter, req *http.Request) {
	logout(res, req)
	tpl.ExecuteTemplate(res, "logout", nil)
	return
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
