package config

type mySQLConfig struct {
	Database           string
	User               string
	Password           string
	Host               string
	Port               int
	Charset            string
	MaxIdleConnections int
	MaxOpenConnections int
	TablePrefix        string
	Url                string
}
