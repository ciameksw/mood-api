package config

import (
	"log"
	"os"
)

type Config struct {
	ServerHost       string
	ServerPort       string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string
	Salt             string
}

func GetConfig() *Config {
	return &Config{
		ServerHost:       getEnv("SERVER_HOST", "localhost"),
		ServerPort:       getEnv("SERVER_PORT", "3001"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDatabase: getEnv("POSTGRES_DATABASE", "auth_db"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		Salt:             getEnv("SALT", "Nb58PsZJlCiO"),
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
