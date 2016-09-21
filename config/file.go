package config

import (
	"io/ioutil"
	"log"

	"github.com/bahusvel/TunnelBeast/auth"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	ListenDev    string
	AuthProvider auth.AuthProvider
}

type configuration struct {
	ListenDev    string
	AuthMethod   string
	AuthProvider map[string]interface{}
}

func LoadConfig(filePath string, conf *Configuration) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed reading config file", err)
	}
	tmpConfig := configuration{}
	err = yaml.Unmarshal(data, &tmpConfig)
	if err != nil {
		log.Fatal("Failed reading config file", err)
	}
	conf.ListenDev = tmpConfig.ListenDev
	switch tmpConfig.AuthMethod {
	case "ldap":
		conf.AuthProvider = auth.LDAPAuth{LDAPAddr: tmpConfig.AuthProvider["ldapaddr"].(string), DCString: tmpConfig.AuthProvider["dcstring"].(string)}
	case "test":
		conf.AuthProvider = auth.TestAuth{Username: tmpConfig.AuthProvider["username"].(string), Password: tmpConfig.AuthProvider["password"].(string)}
	}
}
