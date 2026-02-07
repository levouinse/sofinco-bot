package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken      string
	OwnerID       int64
	OwnerIDs      []int64
	OwnerUsername string
	APIKey        string
	AksesKey      string
	BaseAPIURL    string
}

func Load() *Config {
	ownerID := parseInt64(os.Getenv("OWNER_ID"))
	ownerIDs := []int64{ownerID}
	
	// Parse multiple owner IDs from OWNER_IDS env
	if ownerIDsStr := os.Getenv("OWNER_IDS"); ownerIDsStr != "" {
		for _, idStr := range strings.Split(ownerIDsStr, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
				ownerIDs = append(ownerIDs, id)
			}
		}
	}
	
	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		OwnerID:       ownerID,
		OwnerIDs:      ownerIDs,
		OwnerUsername: os.Getenv("OWNER_USERNAME"),
		APIKey:        os.Getenv("API_KEY"),
		AksesKey:      os.Getenv("AKSES_KEY"),
		BaseAPIURL:    "https://api.betabotz.eu.org",
	}
}

func parseInt64(s string) int64 {
	var result int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int64(c-'0')
		}
	}
	return result
}
