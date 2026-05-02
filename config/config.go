package config

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
