package main

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type LdapConfig struct {
	Url       string
	AdminUser string
	UserAttr  string
	UserOu    string
	LdapDc    string
	LdapPass  string
}

type MailConfig struct {
	Username   string
	Password   string
	SmtpServer string
	SmtpPort   int
}

type Config struct {
	Ldap    *LdapConfig
	Mail    *MailConfig
	Secret  string
	TplPath string
	Tls     bool
	Key     string
	Cert    string
	Port    string
	MaxID   int
	lock    sync.Mutex
}

func validateConfigEntry(entry string, name string) bool {
	if entry == "" {
		log.Printf("Error: %v unset", name)
		return true
	}
	return false
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("guildgate")
	viper.AddConfigPath("/etc/guildgate")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("FATAL: Fatal error reading config file: %v \n", err)
	}
	viper.SetConfigType("yaml")
	c := &Config{}
	l := &LdapConfig{}
	m := &MailConfig{}
	viper.SetDefault("port", "8080")
	viper.SetDefault("secret", "")
	viper.SetDefault("Tls", false)
	//Load configs
	l.Url = viper.GetString("ldapUrl")
	l.AdminUser = viper.GetString("adminUser")
	l.UserAttr = viper.GetString("userAttr")
	l.UserOu = viper.GetString("userOu")
	l.LdapDc = viper.GetString("ldapDc")
	l.LdapPass = viper.GetString("ldapPass")
	c.Secret = viper.GetString("secret")
	c.Tls = viper.GetBool("tls")
	c.Port = viper.GetString("port")
	c.Key = viper.GetString("tls_key")
	c.Cert = viper.GetString("tls_cert")
	c.TplPath = viper.GetString("templates_path")
	m.SmtpServer = viper.GetString("SmtpServer")
	m.Username = viper.GetString("SmtpUsername")
	m.Password = viper.GetString("SmtpPassword")
	m.SmtpPort = viper.GetInt("SmtpPort")

	//Validate configs
	if validateConfigEntry(l.Url, "ldapUrl") || validateConfigEntry(l.AdminUser, "adminUser") || validateConfigEntry(l.UserOu, "userOu") || validateConfigEntry(l.LdapDc, "ldapDc") || validateConfigEntry(l.UserAttr, "userAttr") {
		log.Fatalf("FATAL: Error in config file, bailing")
	}
	c.Ldap = l
	c.Mail = m
	return c, nil
}
