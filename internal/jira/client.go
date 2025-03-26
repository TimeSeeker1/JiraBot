package jira

import (
	"fmt"

	jiralib "github.com/andygrunwald/go-jira"
	"jirabot/internal/config"
)

type Client struct {
	jiraClient *jiralib.Client
}

func NewClient(cfg *config.JiraConfig) (*Client, error) {
	tp := jiralib.BasicAuthTransport{
		Username: cfg.User,
		Password: cfg.Token,
	}

	client, err := jiralib.NewClient(tp.Client(), cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Client{jiraClient: client}, nil
}

func (c *Client) GetIssueService() *jiralib.IssueService {
	return c.jiraClient.Issue
}
