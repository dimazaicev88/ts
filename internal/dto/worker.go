package dto

import "time"

type Worker struct {
	Uid           string    `json:"uid"`
	Name          string    `json:"name"`
	LastConnected time.Time `json:"lastConnected"`
	Status        int8      `json:"status"` // connected, disconnected
}

func (w Worker) IsEmpty() bool {
	return w.Uid == "" &&
		w.Name == "" &&
		w.LastConnected.IsZero() &&
		w.Status == 0
}

type ResponseWorkers struct {
	Items []Worker `json:"items"`
	Total int      `json:"total"`
}
