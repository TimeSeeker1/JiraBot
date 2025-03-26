package config

import "os"

type JiraConfig struct {
	URL     string
	User    string
	Token   string
	Project string
}

type TelegramConfig struct {
	BotToken  string
	ChannelID string
}

func Load() (*JiraConfig, *TelegramConfig) {
	return &JiraConfig{
			URL:     os.Getenv("JIRA_URL"),
			User:    os.Getenv("JIRA_USER"),
			Token:   os.Getenv("JIRA_TOKEN"),
			Project: os.Getenv("JIRA_PROJECT"),
		},
		&TelegramConfig{
			BotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChannelID: os.Getenv("TELEGRAM_CHANNEL_ID"),
		}
}
