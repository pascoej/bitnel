package main

import (
	"reflect"
	"testing"
)

var expectedConfig = &config{
	Database:   "user=andrewtian dbname=bitnel_test sslmode=disable",
	ListenAddr: ":8080",
	BcryptCost: 10,
}

// suppose that config.json.example is working too
func TestLoadConfig(t *testing.T) {
	config, err := loadConfig("config.json.example")
	if err != nil {
		t.Error("can't read config.json.example")
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("config is not equal to expectedConfig")
	}
}
