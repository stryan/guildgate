package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/gomail.v2"
)

func resetLookup(res http.ResponseWriter, req *http.Request) {
	log.Println("POST /reset")
	email := req.FormValue("email")
	uname, err := findLDAPAccountByEmail(email)
	if err != nil {
		log.Printf("Error while looking up account to email password reset to: %v\n. Account may not exist", err)
		http.Redirect(res, req, "/reset/form", 303)
	}
	if uname == "" {
		log.Printf("Error while looking up account to email password reset to: %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Found user %v, generating password token\n", uname)
	token, err := generatePasswordToken(uname)
	fmt.Println(token)
	if err != nil {
		log.Printf("Error generating password token %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	log.Printf("Sending password reset email to %v\n", email)
	go func() {
		err = sendMail(email, uname, token)
		if err != nil {
			log.Printf("Error sending password reset email %v\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	log.Println("Redirecting to next part of password reset")
	http.Redirect(res, req, "/reset/form", 303)
}
func reset(res http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	newPass := req.FormValue("new_password")

	user, err := validateToken(token, true)
	if err != nil {
		log.Printf("Error validing password reset token: %v\n", err)
		http.Redirect(res, req, "/reset/error", 302)
		return
	}
	if user == "" {
		log.Printf("Error resetting password without a username\n")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Attempting to reset password for %v", user)
	err = resetLDAPAccountPassword(user, "", newPass)
	if err == nil {
		log.Printf("reset password for %v\n", user)
		http.Redirect(res, req, "/reset/success", 302)
		return
	}
	log.Printf("failed to reset password for %v:%v\n", user, err)
	http.Redirect(res, req, "/reset/error", 302)
	return

}

func change(res http.ResponseWriter, req *http.Request) {
	oldPass := req.FormValue("old_password")
	newPass := req.FormValue("new_password")

	user := getUserName(req)
	if user == "" {
		log.Printf("Error changing password without a username\n")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Attempting to change password for %v", user)
	err := resetLDAPAccountPassword(user, oldPass, newPass)
	if err != nil {
		log.Printf("failed to change password for %v:%v\n", user, err)
		http.Redirect(res, req, "/change/error", 302)
		return
	}
	log.Printf("change password for %v\n", user)
	http.Redirect(res, req, "/change/success", 302)
	return

}

func sendMail(recp string, uname string, token string) error {
	data := struct {
		Recipient string
		Name      string
		Token     string
		TokenURL  string
	}{
		Recipient: recp,
		Name:      uname,
		Token:     token,
		TokenURL:  url.QueryEscape(token),
	}

	m := gomail.NewMessage()
	m.SetHeader("From", Conf.Mail.Username)
	m.SetHeader("To", recp)
	m.SetHeader("Subject", "Identity Server Password Reset")

	msg := new(bytes.Buffer)

	tpl.ExecuteTemplate(msg, "reset_pass", data)
	m.SetBody("text/plain", string(msg.Bytes()))
	d := gomail.NewDialer(Conf.Mail.SmtpServer, Conf.Mail.SmtpPort, Conf.Mail.Username, Conf.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
