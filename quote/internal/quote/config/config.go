package config

import "github.com/ciameksw/mood-api/pkg/configutil"

type Config struct {
	ServerHost       string
	ServerPort       string
	ExternalQuoteURL string
	RedisAddr        string
}

func GetConfig() *Config {
	return &Config{
		ServerHost:       configutil.GetEnv("SERVER_HOST", "localhost"),
		ServerPort:       configutil.GetEnv("SERVER_PORT", "3002"),
		ExternalQuoteURL: configutil.GetEnv("EXTERNAL_QUOTE_URL", "https://zenquotes.io"),
		RedisAddr:        configutil.GetEnv("REDIS_ADDR", "localhost:6379"),
	}
}
