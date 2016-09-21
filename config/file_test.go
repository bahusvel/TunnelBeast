package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := Configuration{}
	LoadConfig("config.yml", &config)
	fmt.Printf("%+v\n", config)
}
