package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := Configuration{}
	LoadConfig("../test.yml", &config)
	fmt.Printf("%+v\n", config)
}
