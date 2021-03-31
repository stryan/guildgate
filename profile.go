package main

import (
	"log"
	"net/http"
)

func profileEdit(res http.ResponseWriter, req *http.Request) {
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
	dispname := req.FormValue("displayname")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	if dispname != user.DisplayName || firstname != user.FirstName || lastname != user.LastName || email != user.Email {
		log.Printf("updating user %v\n", user.Username)
		user.DisplayName = dispname
		user.FirstName = firstname
		user.LastName = lastname
		user.Email = email
		err = updateLDAPAccountByUser(user)
		if err != nil {
			log.Printf("Error updating user account: %v\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/profile/view", 303)
	return
}
