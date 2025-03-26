package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func New() (*Database, error) {
	db, err := sql.Open("sqlite", "./data.db?_journal=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id TEXT PRIMARY KEY,
            status TEXT,
            priority TEXT,
            transition_time DATETIME,
            notified BOOLEAN
        )`); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetTask(ctx context.Context, id string) (*Task, error) {
	row := d.db.QueryRowContext(ctx,
		"SELECT status, priority, transition_time, notified FROM tasks WHERE id = ?",
		id,
	)

	var t Task
	var ts string
	err := row.Scan(&t.Status, &t.Priority, &ts, &t.Notified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task error: %w", err)
	}

	t.TransitionTime, err = time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, fmt.Errorf("time parse error: %w", err)
	}

	t.ID = id
	return &t, nil
}

func (d *Database) UpsertTask(ctx context.Context, task Task) error {
	_, err := d.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO tasks VALUES (?, ?, ?, ?, ?)`,
		task.ID, task.Status, task.Priority,
		task.TransitionTime.Format(time.RFC3339), task.Notified,
	)
	if err != nil {
		return fmt.Errorf("upsert task failed: %w", err)
	}
	return nil
}

type Task struct {
	ID             string
	Status         string
	Priority       string
	TransitionTime time.Time
	Notified       bool
}
