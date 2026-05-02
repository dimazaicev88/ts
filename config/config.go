package config

import (
	"time"

	"github.com/dimazaicev88/ts/internal/base"
)

type HandlerConfig struct {
	Subject     string
	Concurrency int
	Handler     base.HandlerSubject
	BatchSize   int
	MaxWaitTime time.Duration
	PollingTime time.Duration
}

type DbConfig struct {
	Login           string
	Password        string
	Host            string
	Port            string
	DataBaseName    string
	AllowExecuteSQL bool //Разрешить выполнение sql
	SetTimezone     bool
}

type ServerConfig struct {
	DbConfig       DbConfig
	HttpServerPort int
	SubjectName    string
	NatsURL        string
	StreamName     string
	ConsumerName   string
}

type WorkerConfig struct {
	WorkerName string
	ServerURL  string
}
