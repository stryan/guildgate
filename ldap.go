package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
)

func createLDAPMCAccount(uname, mcuname string) error {
	if uname == "" || mcuname == "" {
		log.Printf("error: missing field\n")
		return errors.New("Missing field")
	}
	url := Conf.Ldap.Url
	newdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, mcuname, Conf.Ldap.MineUserOu, Conf.Ldap.LdapDc)
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	maindn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, uname, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
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
	addReq.Attribute("objectClass", []string{"top", "account"})
	addReq.Attribute("seeAlso", []string{maindn})
	if err := l.Add(addReq); err != nil {
		log.Printf("error adding service:", addReq, err)
		return errors.New("Error creating LDAP account")
	}
	return nil
}

func createLDAPAccount(uname string, pwd string, email string) error {
	if uname == "" || pwd == "" || email == "" {
		log.Printf("error: missing field\n")
		return errors.New("Missing field")
	}
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
	addReq.Attribute("givenName", []string{uname})
	addReq.Attribute("employeeType", []string{"default"})
	addReq.Attribute("employeeNumber", []string{strconv.Itoa(getNextId())})
	addReq.Attribute("displayName", []string{uname})

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
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", Conf.Ldap.UserAttr, uname),
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

func resetLDAPAccountPassword(user string, newPass string) error {
	url := Conf.Ldap.Url
	userdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, user, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
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
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", Conf.Ldap.UserAttr, user),
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
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userdn, "", newPass)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Printf("Password could not be changed: %s", err.Error())
		return errors.New("Error setting password")
	}
	return nil
}

func findLDAPAccountByEmail(email string) (string, error) {
	url := Conf.Ldap.Url
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	basedn := fmt.Sprintf("%v,%v", Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	l, err := ldap.DialURL(url)
	if err != nil {
		return "", err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return "", err
	}
	result, err := l.Search(ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(mail=%s))", email),
		[]string{"dn", Conf.Ldap.UserAttr},
		nil,
	))
	if err != nil {
		return "", err
	}
	if len(result.Entries) != 1 {
		err_text := fmt.Sprintf("Error finding user: Wanted 1 result, got %v\n", len(result.Entries))
		return "", errors.New(err_text)
	}
	entry := result.Entries[0]

	return entry.GetAttributeValue(Conf.Ldap.UserAttr), nil
}
func findLDAPMCAccount(uname string) (string, error) {
	url := Conf.Ldap.Url
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	basedn := fmt.Sprintf("%v,%v", Conf.Ldap.MineUserOu, Conf.Ldap.LdapDc)
	userdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, uname, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	l, err := ldap.DialURL(url)
	if err != nil {
		return "", err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return "", err
	}
	result, err := l.Search(ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=account)(seeAlso=%s))", userdn),
		[]string{"uid"},
		nil,
	))
	if err != nil {
		return "", err
	}
	if len(result.Entries) != 1 {
		err_text := fmt.Sprintf("Error finding user: Wanted 1 result, got %v\n", len(result.Entries))
		return "", errors.New(err_text)
	}
	entry := result.Entries[0]
	return entry.GetAttributeValue("uid"), nil
}
func findLDAPAccountForDisplay(uname string) (User, error) {
	url := Conf.Ldap.Url
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	basedn := fmt.Sprintf("%v,%v", Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
	l, err := ldap.DialURL(url)
	if err != nil {
		return User{}, err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return User{}, err
	}
	result, err := l.Search(ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", Conf.Ldap.UserAttr, uname),
		[]string{"cn", "sn", "givenName", "displayName", "mail", "employeeNumber", "memberOf"},
		nil,
	))
	if err != nil {
		return User{}, err
	}
	if len(result.Entries) != 1 {
		err_text := fmt.Sprintf("Error finding user: Wanted 1 result, got %v\n", len(result.Entries))
		return User{}, errors.New(err_text)
	}
	entry := result.Entries[0]
	groups := entry.GetAttributeValues("memberOf")
	fg := make([]string, 0)
	for _, group := range groups {
		group_s := strings.Split(group, ",")
		group_cn := group_s[0]
		fg = append(fg, strings.Trim(group_cn, "cn="))
	}

	u := User{
		Username:       entry.GetAttributeValue("cn"),
		FirstName:      entry.GetAttributeValue("givenName"),
		LastName:       entry.GetAttributeValue("sn"),
		DisplayName:    entry.GetAttributeValue("displayName"),
		Email:          entry.GetAttributeValue("mail"),
		ID:             entry.GetAttributeValue("employeeNumber"),
		Groups:         groups,
		FriendlyGroups: fg,
	}
	return u, nil
}
func updateLDAPAccountByUser(user User) error {
	url := Conf.Ldap.Url
	userdn := fmt.Sprintf("%v=%v,%v,%v", Conf.Ldap.UserAttr, user.Username, Conf.Ldap.UserOu, Conf.Ldap.LdapDc)
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
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", Conf.Ldap.UserAttr, user.Username),
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
	modify := ldap.NewModifyRequest(userdn, nil)
	modify.Replace("mail", []string{user.Email})
	modify.Replace("givenName", []string{user.FirstName})
	modify.Replace("sn", []string{user.LastName})
	modify.Replace("displayName", []string{user.DisplayName})
	err = l.Modify(modify)

	if err != nil {
		return err
	}
	return nil
}
func findLDAPMaxID() (int, error) {
	url := Conf.Ldap.Url
	binddn := fmt.Sprintf("%v,%v", Conf.Ldap.AdminUser, Conf.Ldap.LdapDc)
	basedn := fmt.Sprintf("%v,%v", Conf.Ldap.UserOu, Conf.Ldap.LdapDc)

	l, err := ldap.DialURL(url)
	if err != nil {
		return -1, err
	}
	defer l.Close()
	err = l.Bind(binddn, Conf.Ldap.LdapPass)
	if err != nil {
		return -1, err
	}
	result, err := l.Search(ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(employeeNumber=*))"),
		[]string{"employeeNumber"},
		nil,
	))
	if err != nil {
		return -1, err
	}
	maxId := 0
	for _, entry := range result.Entries {
		i, err := strconv.Atoi(entry.GetAttributeValue("employeeNumber"))
		if err != nil {
			return -1, err
		}
		if i > maxId {
			maxId = i
		}
	}
	return maxId + 1, nil

}

func getNextId() int {
	if Conf.MaxID == 0 {
		return -1
	}
	Conf.lock.Lock()
	i := Conf.MaxID
	Conf.MaxID = Conf.MaxID + 1
	Conf.lock.Unlock()
	return i
}
