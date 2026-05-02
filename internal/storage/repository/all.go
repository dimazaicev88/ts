package repository

import "github.com/uptrace/bun"

type AllRepository struct {
	Task   ITask
	Worker IWorker
}

func NewAll(db *bun.DB) AllRepository {
	return AllRepository{
		Task:   NewTask(db),
		Worker: NewWorker(db),
	}
}
