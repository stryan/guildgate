package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Url       string
	AdminUser string
	UserAttr  string
	UserOu    string
	LdapDc    string
	Secret    string
}

func validateConfigEntry(entry string, name string) bool {
	if entry == "" {
		log.Printf("Error: %v unset", name)
		return false
	}
	return true
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("guildgate.yaml")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("FATAL: Fatal error reading config file: %v \n", err)
	}
	viper.SetConfigType("yaml")
	c := &Config{}
	//Load configs
	c.Url = viper.GetString("ldapUrl")
	c.AdminUser = viper.GetString("adminUser")
	c.UserAttr = viper.GetString("userAttr")
	c.UserOu = viper.GetString("userOu")
	c.LdapDc = viper.GetString("ldapDc")
	c.Secret = viper.GetString("secret")

	//Validate configs
	if validateConfigEntry(c.Url, "ldapUrl") || validateConfigEntry(c.AdminUser, "adminUser") || validateConfigEntry(c.UserOu, "userOu") || validateConfigEntry(c.LdapDc, "ldapDc") || validateConfigEntry(c.UserAttr, "userAttr") {
		log.Fatalf("FATAL: Error in config file, bailing")
	}

	return c, nil
}
