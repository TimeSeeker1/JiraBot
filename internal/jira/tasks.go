package jira

import (
	"context"
	"fmt"
	"time"

	jiralib "github.com/andygrunwald/go-jira"
)

type TaskService struct {
	issueService *jiralib.IssueService
}

func NewTaskService(client *Client) *TaskService {
	return &TaskService{
		issueService: client.GetIssueService(),
	}
}

func (s *TaskService) GetProjectTasks(ctx context.Context, project string) ([]Task, error) {
	jql := fmt.Sprintf("project = %s AND status NOT IN (Closed, Resolved)", project)

	issues, _, err := s.issueService.SearchWithContext(ctx, jql, &jiralib.SearchOptions{
		MaxResults: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	tasks := make([]Task, len(issues))
	for i, issue := range issues {
		tasks[i] = convertIssueToTask(issue) // Убрано разыменование (*)
	}
	return tasks, nil
}

func (s *TaskService) GetTaskHistory(taskID string) ([]StatusTransition, error) {
	issue, _, err := s.issueService.Get(taskID, &jiralib.GetQueryOptions{
		Expand: "changelog",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// Проверка на nil для Changelog
	if issue.Changelog == nil {
		return nil, nil
	}

	transitions := make([]StatusTransition, 0)
	for _, history := range issue.Changelog.Histories {
		// Проверка на nil для Items
		if history.Items == nil {
			continue
		}

		for _, item := range history.Items {
			if item.Field == "status" {
				// Парсинг времени через стандартный time.Parse
				transitionTime, err := time.Parse(time.RFC3339, history.Created)
				if err != nil {
					return nil, fmt.Errorf("time parse error: %w", err)
				}

				transitions = append(transitions, StatusTransition{
					FromStatus: item.FromString,
					ToStatus:   item.ToString,
					Date:       transitionTime,
				})
			}
		}
	}

	return transitions, nil
}
