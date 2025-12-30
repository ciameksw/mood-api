package config

import "github.com/ciameksw/mood-api/pkg/configutil"

type Config struct {
	ServerHost       string
	ServerPort       string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string
}

func GetConfig() *Config {
	return &Config{
		ServerHost:       configutil.GetEnv("SERVER_HOST", "localhost"),
		ServerPort:       configutil.GetEnv("SERVER_PORT", "3002"),
		PostgresHost:     configutil.GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     configutil.GetEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     configutil.GetEnv("POSTGRES_USER", "user"),
		PostgresPassword: configutil.GetEnv("POSTGRES_PASSWORD", "password"),
		PostgresDatabase: configutil.GetEnv("POSTGRES_DATABASE", "mood_api_db"),
		PostgresSSLMode:  configutil.GetEnv("POSTGRES_SSLMODE", "disable"),
	}
}
