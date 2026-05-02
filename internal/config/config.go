package config

import "github.com/dimazaicev88/ts/internal/storage"

type Config struct {
	DbConfig       storage.DbConfig
	HttpServerPort int
	SubjectName    string
	NatsURL        string
	StreamName     string
	ConsumerName   string
}
