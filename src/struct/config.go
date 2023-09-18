package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
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
