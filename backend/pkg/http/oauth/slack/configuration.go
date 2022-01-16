package slack

import (
	"log"
	"os"
)

type Configuration struct {
	ClientID     string
	ClientSecret string
}

func (c *Configuration) Parse() error {
	clientID := os.Getenv("SLACK_CLIENT_ID")
	if clientID == "" {
		log.Printf("[WARN] SLACK_CLIENT_ID is not set")
	}
	c.ClientID = clientID

	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	if clientSecret == "" {
		log.Printf("[WARN] SLACK_CLIENT_SECRET is not set")
	}
	c.ClientSecret = clientSecret

	return nil
}
