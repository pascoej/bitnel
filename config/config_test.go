package config

import (
	"reflect"
	"testing"
)

var expectedConfig = &Config{
	Database:   "user=andrewtian dbname=bitnel_test sslmode=disable",
	ListenAddr: ":8080",
	BcryptCost: 10,
}

// suppose that config.json.example is working too
func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("test/example.json")
	if err != nil {
		t.Error("can't read test/example.json")
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("config is not equal to expectedConfig")
	}
}
