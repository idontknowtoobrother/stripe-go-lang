package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return os.Getenv(key)
}
