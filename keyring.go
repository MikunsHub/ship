package main

import (
	"os"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "ship-cli"
	accountName = "gemini-api-key"
)

func SaveAPIKey(key string) error {
	return keyring.Set(serviceName, accountName, key)
}

func GetAPIKey() string {
	if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" {
		return envKey
	}

	key, err := keyring.Get(serviceName, accountName)
	if err == nil {
		return key
	}

	return ""
}

func DeleteAPIKey() error {
	return keyring.Delete(serviceName, accountName)
}

func HasAPIKey() bool {
	return GetAPIKey() != ""
}
