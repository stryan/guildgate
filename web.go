package main

import (
	"log"
	"net/http"
)

func profilePage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /profile")
	uname := getUserName(req)
	if uname == "" {
		http.Redirect(res, req, "/", 302)
	}
	user, err := findLDAPAccountForDisplay(uname)
	if err != nil {
		log.Printf("Error loading profile: %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data := struct {
		Title    string
		Username string
		LoggedIn bool
		User     User
	}{
		"Profile",
		uname,
		true,
		user,
	}
	tpl.ExecuteTemplate(res, "profile", data)
}

func resetPageFront(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /passwordreset")
	u := getUserName(req)
	if u != "" {
		http.Redirect(res, req, "/", 302) //TODO create password change form, direct to that
	} else {
		data := struct {
			Title    string
			Username string
			LoggedIn bool
		}{
			"Reset Password",
			"Unregistered",
			false,
		}
		tpl.ExecuteTemplate(res, "reset_password_page_front", data)
	}
}

func resetPageBack(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /passwordresetform")
	u := getUserName(req)
	if u != "" {
		http.Redirect(res, req, "/", 302) //TODO create password change form, direct to that
	} else {
		data := struct {
			Title    string
			Username string
			LoggedIn bool
		}{
			"Reset Password",
			"Unregistered",
			false,
		}
		tpl.ExecuteTemplate(res, "reset_password_page_back", data)
	}
}
func resetSuccessPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /resetsuccess")
	data := struct {
		Title    string
		Username string
		LoggedIn bool
	}{
		"Reset Password Success",
		"Unregistered",
		false,
	}
	tpl.ExecuteTemplate(res, "reset_success", data)
	return
}
func resetErrorPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /reseterror")
	data := struct {
		Title    string
		Username string
		LoggedIn bool
		Error    string
	}{
		"Reset Password Failure",
		"Unregistered",
		false,
		"Undefined",
	}
	tpl.ExecuteTemplate(res, "reset_error", data)
	return
}
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
	log.Println("GET /logout")
	logout(res, req)
	tpl.ExecuteTemplate(res, "logout", nil)
	return
}

func tokenPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /token")
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
