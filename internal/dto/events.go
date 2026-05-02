package dto

import (
	"time"

	"github.com/segmentio/encoding/json"
)

type (
	EventData struct {
		EventName string        `json:"eventName"`
		Payload   json.RawValue `json:"payload"`
	}

	EventWorkerConnect struct {
		WorkerName string `json:"workerName"`
		Hostname   string `json:"hostname"`
	}

	EventAnswerWorkerConnect struct {
		WorkerUid  string `json:"workerUid"`
		WorkerName string `json:"workerName"`
	}

	EventNewTask struct {
		Uid        string    `json:"uid"`
		Status     string    `json:"status"`
		DateCreate time.Time `json:"dateCreate"`
		DateUpdate time.Time `json:"dateUpdate"`
		Retention  uint64    `json:"retention"`
		ErrorMsg   string    `json:"errorMsg"`
		Subject    string    `json:"subject"`
		WorkerUid  string    `json:"workerUid"`
		Timeout    uint64    `json:"timeout"`
		Payload    string    `json:"payload"`
	}

	EventUpdateTask struct {
		Uid    string            `json:"uid"`
		Fields map[string]string `json:"fields"`
	}
)
