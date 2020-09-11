package config

type mySQLConfig struct {
	Database           string `json:"database"`
	User               string `json:"user"`
	Password           string `json:"password"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	Charset            string `json:"charset"`
	TablePrefix        string `json:"table_prefix"`
	MaxIdleConnections int    `json:"max_idle_connections"`
	MaxOpenConnections int    `json:"max_open_connections"`
	Url                string `json:"-"`
}
