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
	apiKey := os.Getenv("API_KEY")
	return apiKey, nil
}
