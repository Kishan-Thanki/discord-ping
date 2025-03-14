package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

var (
	Token     string
	BotPrefix string

	config *configStruct
)

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
