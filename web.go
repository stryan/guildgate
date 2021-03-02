package main

import (
	"log"
	"net/http"
	"net/url"
)

func profilePage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /profile/view")
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
func profileEditPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /profile/edit")
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
	tpl.ExecuteTemplate(res, "profile_edit", data)
}

func minecraftPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /minecraft")
	u := getUserName(req)
	if u == "" {
		http.Redirect(res, req, "/", 302)
	} else {
		mcuser, err := findLDAPMCAccount(u)
		mclink := true
		if err != nil {
			mclink = false
			mcuser = "N/A"
		}
		data := struct {
			Title            string
			Username         string
			LoggedIn         bool
			Linked           bool
			MinecraftAccount string
		}{
			"Link Minecraft Account",
			u,
			true,
			mclink,
			mcuser,
		}
		tpl.ExecuteTemplate(res, "minecraft", data)
	}
}
func minecraftLinkSuccessPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /minecraft/link/success")
	u := getUserName(req)
	if u == "" {
		http.Redirect(res, req, "/404", 302)
	}
	genericSuccessPage(res, "Minecraft Link Success", u, true, "Succesfully linked Minecraft account.")
	return
}
func minecraftLinkErrorPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /minecraft/link/error")
	u := getUserName(req)
	if u == "" {
		http.Redirect(res, req, "/404", 302)
	}
	genericErrorPage(res, "Minecraft Link Failure", u, true, "Undefined", "link Minecraft account.")
	return
}
func resetPageFront(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /reset")
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
	log.Println("GET /reset/form")
	u := getUserName(req)
	token := ""
	if u != "" {
		http.Redirect(res, req, "/", 302) //TODO create password change form, direct to that
	} else {
		keys, ok := req.URL.Query()["token"]
		if !ok || len(keys[0]) < 1 {
			token = ""
		} else {
			token = keys[0]
		}
		data := struct {
			Title    string
			Username string
			LoggedIn bool
			Token    string
		}{
			"Reset Password",
			"Unregistered",
			false,
			token,
		}
		tpl.ExecuteTemplate(res, "reset_password_page_back", data)
	}
}
func resetSuccessPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /reset/success")
	genericSuccessPage(res, "Reset Password Success", "Unregistered", false, "Succesfully Reset Password")
	return
}
func resetErrorPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /reset/error")
	genericErrorPage(res, "Reset Password Failure", "Unregistered", false, "Undefined", "reset password")
	return
}
func signupPage(res http.ResponseWriter, req *http.Request) {
	log.Println("GET /register")
	u := getUserName(req)
	secret := ""
	if u != "" {
		http.Redirect(res, req, "/", 302)
	} else {
		keys, ok := req.URL.Query()["secret"]
		if !ok || len(keys[0]) < 1 {
			secret = ""
		} else {
			secret = keys[0]
		}
		data := struct {
			Title    string
			Username string
			LoggedIn bool
			Secret   string
		}{
			"Register",
			"Unregistered",
			false,
			secret,
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
		TokenURL string
	}{
		"Token Generation",
		u,
		true,
		token,
		url.QueryEscape(token),
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

func genericSuccessPage(res http.ResponseWriter, title string, uname string, login bool, action string) {
	data := struct {
		Title    string
		Username string
		LoggedIn bool
		Action   string
	}{
		title,
		uname,
		login,
		action,
	}
	tpl.ExecuteTemplate(res, "generic_success", data)
	return
}
func genericErrorPage(res http.ResponseWriter, title string, uname string, login bool, err string, action string) {
	data := struct {
		Title    string
		Username string
		LoggedIn bool
		Error    string
		Action   string
	}{
		title,
		uname,
		login,
		err,
		action,
	}
	tpl.ExecuteTemplate(res, "generic_error", data)
	return
}
