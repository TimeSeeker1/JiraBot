package monitoring

import (
	"context"
	"log"
	"time"

	"jirabot/internal/jira"
	"jirabot/internal/storage"
	"jirabot/internal/telegram"
)

type Service struct {
	jiraService   *jira.TaskService
	tgBot         *telegram.Bot
	storage       *storage.Database
	checkInterval time.Duration
	project       string
}

func NewService(
	jiraService *jira.TaskService,
	tgBot *telegram.Bot,
	db *storage.Database,
	interval time.Duration,
	project string,
) *Service {
	return &Service{
		jiraService:   jiraService,
		tgBot:         tgBot,
		storage:       db,
		checkInterval: interval,
		project:       project,
	}
}

func (s *Service) Start() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.checkTasks()
	}
}

func (s *Service) checkTasks() {
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	tasks, err := s.jiraService.GetProjectTasks(ctx, s.project)
	if err != nil {
		log.Printf("Get tasks error: %v", err)
		return
	}

	for _, task := range tasks {
		storedTask, err := s.storage.GetTask(ctx, task.ID)
		if err != nil {
			log.Printf("Get task error: %v", err)
			continue
		}

		history, err := s.jiraService.GetTaskHistory(task.ID)
		if err != nil {
			log.Printf("Get history error: %v", err)
			continue
		}

		if len(history) == 0 {
			continue
		}

		lastTransition := history[len(history)-1].Date
		limit, ok := getTimeLimit(task.Status, task.Priority)
		if !ok {
			continue
		}

		elapsed := time.Since(lastTransition)
		if elapsed > limit && (storedTask == nil || !storedTask.Notified) {
			if err := s.tgBot.SendAlert(task.Key, task.Status, task.Priority, limit); err != nil {
				log.Printf("Send alert error: %v", err)
				continue
			}

			if err := s.storage.UpsertTask(ctx, storage.Task{
				ID:             task.ID,
				Status:         task.Status,
				Priority:       task.Priority,
				TransitionTime: lastTransition,
				Notified:       true,
			}); err != nil {
				log.Printf("Update task error: %v", err)
			}
		}
	}
}

func getTimeLimit(status, priority string) (time.Duration, bool) {
	switch status {
	case "Open":
		return 15 * time.Minute, true
	case "In Progress":
		switch priority {
		case "Medium", "High":
			return 1 * time.Hour, true
		default:
			return 15 * time.Minute, true
		}
	case "Hold":
		return 6 * time.Hour, true
	case "Need Info", "L2 Escalation out":
		switch priority {
		case "Medium", "High":
			return 3 * time.Hour, true
		default:
			return 45 * time.Minute, true
		}
	case "Prod fixed":
		return 30 * time.Minute, true
	default:
		return 0, false
	}
}
