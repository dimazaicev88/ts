package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Task Таблица хранит информацию о задачах
type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t_tasks"`
	Uid           string    `bun:"uid"`
	Status        string    `bun:"status"`
	DateCreate    time.Time `bun:"date_create"`
	DateUpdate    time.Time `bun:"date_update"`
	Result        string    `bun:"result"`  // Результат работы воркера
	Subject       string    `bun:"subject"` // Nats subject
	WorkerUid     string    `bun:"worker_uid"`
	Timeout       uint64    `bun:"timeout"`
	Retention     uint64    `bun:"retention"`
	CompletedAt   uint64    `bun:"completed_at"`
	Payload       string    `bun:"payload"`
	ErrorMsg      string    `bun:"error_msg"`
}

// Worker Таблица хранит информацию о воркерах
type Worker struct {
	bun.BaseModel `bun:"table:workers,alias:t_workers"`
	Uid           string    `bun:"uid"`
	Name          string    `bun:"name"`
	LastConnected time.Time `bun:"last_connected"`
	Status        int8      `bun:"status"` // connected, disconnected
}
