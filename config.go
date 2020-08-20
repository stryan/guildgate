package main

import (
	"log"

	"github.com/spf13/viper"
)

type LdapConfig struct {
	Url       string
	AdminUser string
	UserAttr  string
	UserOu    string
	LdapDc    string
}

type Config struct {
	Ldap   *LdapConfig
	Secret string
	Tls    bool
	Port   string
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
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("FATAL: Fatal error reading config file: %v \n", err)
	}
	viper.SetConfigType("yaml")
	c := &Config{}
	l := &LdapConfig{}
	viper.SetDefault("port", "8080")
	//Load configs
	l.Url = viper.GetString("ldapUrl")
	l.AdminUser = viper.GetString("adminUser")
	l.UserAttr = viper.GetString("userAttr")
	l.UserOu = viper.GetString("userOu")
	l.LdapDc = viper.GetString("ldapDc")
	c.Secret = viper.GetString("secret")
	c.Tls = viper.GetBool("tls")
	c.Port = viper.GetString("port")

	//Validate configs
	if validateConfigEntry(l.Url, "ldapUrl") || validateConfigEntry(l.AdminUser, "adminUser") || validateConfigEntry(l.UserOu, "userOu") || validateConfigEntry(l.LdapDc, "ldapDc") || validateConfigEntry(l.UserAttr, "userAttr") {
		log.Fatalf("FATAL: Error in config file, bailing")
	}
	c.Ldap = l
	return c, nil
}
