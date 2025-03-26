package jira

import (
	"time"

	jiralib "github.com/andygrunwald/go-jira"
)

type Task struct {
	ID             string
	Key            string
	Status         string
	Priority       string
	TransitionTime time.Time
}

type StatusTransition struct {
	FromStatus string
	ToStatus   string
	Date       time.Time
}

func convertIssueToTask(issue jiralib.Issue) Task {
	task := Task{
		ID:  issue.ID,
		Key: issue.Key,
	}

	if issue.Fields != nil {
		if issue.Fields.Status != nil {
			task.Status = issue.Fields.Status.Name
		}
		if issue.Fields.Priority != nil {
			task.Priority = issue.Fields.Priority.Name
		}
	}

	return task
}
