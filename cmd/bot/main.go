package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jirabot/internal/config"
	"jirabot/internal/jira"
	"jirabot/internal/monitoring"
	"jirabot/internal/storage"
	"jirabot/internal/telegram"
)

func main() {
	// Load configuration
	jiraCfg, tgCfg := config.Load()

	// Initialize Jira client
	jiraClient, err := jira.NewClient(jiraCfg)
	if err != nil {
		log.Fatalf("Jira client error: %v", err)
	}
	taskService := jira.NewTaskService(jiraClient)

	// Initialize Telegram bot
	tgBot, err := telegram.NewBot(tgCfg)
	if err != nil {
		log.Fatalf("Telegram bot error: %v", err)
	}

	// Initialize storage
	db, err := storage.New()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Create monitoring service
	monitor := monitoring.NewService(
		taskService,
		tgBot,
		db,
		1*time.Minute,
		jiraCfg.Project,
	)

	// Handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Service started. Press Ctrl+C to stop.")
	go monitor.Start()

	<-stop
	log.Println("Shutting down...")
}
