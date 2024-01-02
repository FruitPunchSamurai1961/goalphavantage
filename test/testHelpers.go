package test

import (
	"github.com/joho/godotenv"
	"os"
)

func getApiKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	apiKey, _ := os.LookupEnv("API_KEY")
	return apiKey, nil
}
