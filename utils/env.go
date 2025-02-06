package utils

import (
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

var (
	onceEnv sync.Once
)

// LoadEnv ensures that the .env file is only loaded once
func LoadEnv() {
	var err error
	onceEnv.Do(func() {
		err = godotenv.Load()
		if err != nil {
			log.Warn(".env file not found, using defaults.")
		}
	})
}

func AllowedList() []string {
	allowedListStr := os.Getenv("APP_ALLOWED_REDIRECT")
	return strings.Split(allowedListStr, ",")
}

func IsAllowedRedirect(allowed string) bool {
	if allowed == "" {
		return true
	}

	allowedList := AllowedList()
	for _, v := range allowedList {
		if v == allowed {
			return true
		}
	}
	return false
}
