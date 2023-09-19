package config

import (
	"os"

	logger "github.com/Malakhov-cmd/GPT-Telegram-Go.git/src/util"
	"gopkg.in/yaml.v3"
)

var (
	log = logger.GetLogger()
)

type Config struct {
	API_Keys struct {
		Telegram_Keys []string `yaml:"Telegram_Keys"`
		Openai_Keys   []string `yaml:"Openai_Keys"`
	} `yaml:"API_Keys"`
}

func GetConfig() Config {
	f, err := os.ReadFile("../config.yml")
	if err != nil {
		log.Fatal("Не удалось прочитать файл конфигурации")
	}

	var config Config
	err = yaml.Unmarshal(f, &config)

	if err != nil {
		log.Fatal("Не удалось десерриализовать объект конфигурации")
	}

	return config
}
