package config

import "github.com/ciameksw/mood-api/pkg/configutil"

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
		ServerHost: configutil.GetEnv("SERVER_HOST", "localhost"),
		ServerPort: configutil.GetEnv("SERVER_PORT", "3002"),
		AdviceURL:  configutil.GetEnv("ADVICE_URL", "http://localhost:3003"),
		MoodURL:    configutil.GetEnv("MOOD_URL", "http://localhost:3002"),
		AuthURL:    configutil.GetEnv("AUTH_URL", "http://localhost:3001"),
		QuoteURL:   configutil.GetEnv("QUOTE_URL", "http://localhost:3004"),
	}
}
