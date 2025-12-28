package config

import (
	"log"
	"os"
)

type Config struct {
	ServerHost string
	ServerPort string
	AdviceURL  string
	MoodURL    string
	AuthURL    string
	QuoteURL   string
}

func GetConfig() *Config {
	return &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "3002"),
		AdviceURL:  getEnv("ADVICE_URL", "http://localhost:3003"),
		MoodURL:    getEnv("MOOD_URL", "http://localhost:3002"),
		AuthURL:    getEnv("AUTH_URL", "http://localhost:3001"),
		QuoteURL:   getEnv("QUOTE_URL", "http://localhost:3004"),
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
