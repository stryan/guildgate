package main

import (
	"log"
	"net/http"
)

func minecraftLink(res http.ResponseWriter, req *http.Request) {
	uname := getUserName(req)
	if uname == "" {
		http.Redirect(res, req, "/", 302)
	}
	mcname := req.FormValue("mcusername")
	if mcname != "" {
		log.Printf("linked MC %v to LDAP %v\n", mcname, uname)
		err := createLDAPMCAccount(uname, mcname)
		if err != nil {
			log.Printf("Error linking MC account: %v\n", err)
			http.Redirect(res, req, "/minecraft/link/error", 302)
		} else {
			http.Redirect(res, req, "/minecraft/link/success", 302)
		}
	} else {
		log.Println("couldn't get MC username")
		http.Redirect(res, req, "/minecraft/link/error", 302)
	}
	return
}
