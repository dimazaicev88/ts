package storage

type DbConfig struct {
	Login           string
	Password        string
	Host            string
	Port            string
	DataBaseName    string
	AllowExecuteSQL bool //Разрешить выполнение sql
	SetTimezone     bool
}
