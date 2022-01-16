package slack

import (
	"log"
	"os"
)

type Configuration struct {
	SigningSecret  string
	BotAccessToken string
}

func (c *Configuration) Parse() error {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		log.Printf("[WARN] SLACK_SIGNING_SECRET is not set, webhooks' signature will not be verified")
	}
	c.SigningSecret = signingSecret

	slackBotAccessToken := os.Getenv("SLACK_BOT_ACCESS_TOKEN")
	if slackBotAccessToken == "" {
		log.Printf("[WARN] SLACK_BOT_ACCESS_TOKEN is not set, slack notifications won't work")
	}
	c.BotAccessToken = slackBotAccessToken

	return nil
}
