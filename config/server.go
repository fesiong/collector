package config

type serverConfig struct {
	SiteName string `json:"site_name"`
	Host     string `json:"host"`
	Env      string `json:"env"`
	LogLevel string `json:"log_level"`
	Port     int    `json:"port"`
	ExecPath string `json:"-"`
}
