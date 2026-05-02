package dto

import "time"

type ResponseTask struct {
	Uid         string    `json:"uid"`
	Payload     string    `json:"payload"`
	Status      string    `json:"status"`
	DateCreate  time.Time `json:"dateCreate"`
	DateUpdate  time.Time `json:"dateUpdate"`
	Result      string    `json:"result"`  // Результат работы воркера
	Subject     string    `json:"subject"` // Nats subject
	WorkerUid   string    `json:"workerUid"`
	Timeout     uint64    `json:"timeout"`
	Retention   uint64    `json:"retention"`
	CompletedAt uint64    `json:"completedAt"`
	ErrorMsg    string    `json:"errorMsg"`
}

type ResponseTasks struct {
	Items []ResponseTask `json:"items"`
	Total int            `json:"total"`
}
