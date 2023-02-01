package config

import (
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := Configuration{}
	LoadConfig("../test.yml", &config)
	log.Printf("%+v\n", config)
}
