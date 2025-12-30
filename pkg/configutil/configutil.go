package configutil

import (
	"log"
	"os"
)

func GetEnv(key, df string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Using default value for %s (%s)", key, df)
		return df
	}
	return val
}
