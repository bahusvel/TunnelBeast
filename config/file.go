package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/bahusvel/TunnelBeast/auth"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	DBpath       string
	Https        string
	Path         string
	Ports        []string
	Domainname   string
	AuthProvider auth.AuthProvider
}

type configuration struct {
	DBpath       string
	Https        string
	Path         string
	AuthMethod   string
	Ports        []string
	Domainname   string
	AuthProvider map[string]interface{}
}

func setField(obj interface{}, name string, value interface{}) error {
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

func fillStruct(target interface{}, data map[string]interface{}) error {
	for k, v := range data {
		err := setField(target, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsePorts(ports []string) []string {
	set := map[int]struct{}{}
	for _, port := range ports {
		switch strings.Count(port, "-") {
		case 0:
			p, err := strconv.Atoi(port)
			if err != nil {
				goto err
			}
			if p < 0 || p > 65536 {
				goto err
			}
			set[p] = struct{}{}
		case 1:
			parts := strings.Split(port, "-")
			from, errFrom := strconv.Atoi(parts[0])
			to, errTo := strconv.Atoi(parts[1])
			if errFrom != nil || errTo != nil {
				goto err
			}
			if from < 0 || from > 65536 || to < 0 || to >= 65536 || from > to {
				goto err
			}
			for i := from; i <= to; i++ {
				set[i] = struct{}{}
			}
		default:
			goto err
		}
		continue
	err:
		log.Fatal("Port definition is invalid ", port)
	}
	portInts := []int{}
	for port := range set {
		portInts = append(portInts, port)
	}
	sort.Ints(portInts)
	portStrings := []string{}
	for _, port := range portInts {
		portStrings = append(portStrings, strconv.Itoa(port))
	}
	return portStrings
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
	conf.DBpath = tmpConfig.DBpath
	conf.Https = tmpConfig.Https
	conf.Path = tmpConfig.Path
	conf.Ports = parsePorts(tmpConfig.Ports)
	conf.Domainname = tmpConfig.Domainname
	switch tmpConfig.AuthMethod {
	case "ldap":
		tmpProvider := auth.LDAPAuth{}
		fillStruct(&tmpProvider, tmpConfig.AuthProvider)
		conf.AuthProvider = tmpProvider
	case "test":
		tmpProvider := auth.TestAuth{}
		fillStruct(&tmpProvider, tmpConfig.AuthProvider)
		conf.AuthProvider = tmpProvider
	case "local":
		tmpProvider := auth.LocalAuth{}
		fillStruct(&tmpProvider, tmpConfig.AuthProvider)
		conf.AuthProvider = tmpProvider
	}
}
