package helper

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func GetEnv(key string) (string, error) {
	env, err := os.LookupEnv(key)
	if !err || env == "" {
		return env, errors.New("environment variable '" + key + "' empty or not found")
	}
	return env, nil
}

func init() {
	godotenv.Load()
}
