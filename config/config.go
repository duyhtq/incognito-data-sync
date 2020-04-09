package config

import (
	"encoding/json"
	"log"
	"os"
)

var config *Config

func init() {
	file, err := os.Open("config/conf.json")
	if err != nil {
		log.Println("error:", err)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	v := Config{}
	err = decoder.Decode(&v)
	if err != nil {
		log.Println("error:", err)
		panic(err)
	}
	config = &v
}

func GetConfig() *Config {
	return config
}

type IncognitoConfig struct {
	ChainEndpoint   string `json:"chain_endpoint"`
	NetWorkEndPoint string `json:"network_endpoint"`
	ApiEndPoint     string `json:"api_endpoint"`
}

type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	BaseURL string `json:"base_url"`
	Db      string `json:"db"`

	SentryDSN string `json:"sentry_dsn"`

	// 0. Incognito chain config:
	Incognito IncognitoConfig `json:"incognito"`
}
