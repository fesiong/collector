package main

import (
	"collector/app/bootstrap"
	"collector/config"
)

func main() {
	b := bootstrap.New(config.ServerConfig.Port, config.ServerConfig.LogLevel)
	b.Serve()
}