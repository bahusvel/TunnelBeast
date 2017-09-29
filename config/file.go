package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/bahusvel/TunnelBeast/auth"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	ListenDev    string
	Ports        []string
    Domainname   string
	AuthProvider auth.AuthProvider
}

type configuration struct {
	ListenDev    string
	AuthMethod   string
	Ports        []string
    Domainname   string
	AuthProvider map[string]interface{}
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(target interface{}, data map[string]interface{}) error {
	for k, v := range data {
		err := SetField(target, k, v)
		if err != nil {
			return err
		}
	}
	return nil
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
	conf.Ports = tmpConfig.Ports
    conf.Domainname = tmpConfig.Domainname
	switch tmpConfig.AuthMethod {
	case "ldap":
		tmpProvider := auth.LDAPAuth{}
		FillStruct(&tmpProvider, tmpConfig.AuthProvider)
		conf.AuthProvider = tmpProvider
	case "test":
		tmpProvider := auth.TestAuth{}
		FillStruct(&tmpProvider, tmpConfig.AuthProvider)
		conf.AuthProvider = tmpProvider
	}
}
