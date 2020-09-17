package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

func createLDAPAccount(uname string, pwd string, email string) error {
	url := Conf.Ldap.Url
	newdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, uname, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	l, err := ldap.DialURL(url)
	if err != nil {
		return err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return err
	}
	addReq := ldap.NewAddRequest(newdn, []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "inetOrgPerson"})
	addReq.Attribute("cn", []string{uname})
	addReq.Attribute("mail", []string{email})
	addReq.Attribute("sn", []string{"The Nameless"})

	if err := l.Add(addReq); err != nil {
		log.Printf("error adding service:", addReq, err)
		return errors.New("Error creating LDAP account")
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest(newdn, "", pwd)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Printf("Password could not be changed: %s", err.Error())
		return errors.New("Error setting password")
	}
	return nil
}

func loginLDAPAccount(uname string, pwd string) error {
	url := Conf.Ldap.Url
	userdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, uname, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	basedn := fmt.Sprintf("%v,%v", Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	l, err := ldap.DialURL(url)
	if err != nil {
		return err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return err
	}
	result, err := l.Search(ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", uname),
		[]string{"dn"},
		nil,
	))
	if err != nil {
		return err
	}
	if len(result.Entries) != 1 {
		err_text := fmt.Sprintf("Error finding login user: Wanted 1 result, got %v\n", len(result.Entries))
		return errors.New(err_text)
	}
	err = l.Bind(userdn, pwd)
	if err != nil {
		return err
	}
	return nil
}
