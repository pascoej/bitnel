package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Fields are uppercase so that they can be exported (accessed) by the json pkg
type config struct {
	Database   string `json:"database"`
	ListenAddr string `json:"listenAddr"`
	BcryptCost int    `json:"bcryptCost"`
}

func loadConfig(filename string) (*config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("cannot read config file")
	}

	var m config
	err = json.Unmarshal(b, &m)

	if err != nil {
		return nil, errors.New("cannot unmarshal json")
	}

	return &m, nil
}
