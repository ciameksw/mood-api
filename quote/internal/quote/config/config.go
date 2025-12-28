package config

import (
	"log"
	"os"
)

type Config struct {
	ServerHost       string
	ServerPort       string
	ExternalQuoteURL string
}

func GetConfig() *Config {
	return &Config{
		ServerHost:       getEnv("SERVER_HOST", "localhost"),
		ServerPort:       getEnv("SERVER_PORT", "3002"),
		ExternalQuoteURL: getEnv("EXTERNAL_QUOTE_URL", "https://zenquotes.io"),
	}
}

func getEnv(key, df string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Using default value for %s (%s)", key, df)
		return df
	}
	return val
}
