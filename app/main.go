package main

import (
	"collector"
	"collector/app/bootstrap"
	"collector/config"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	go listenSignal()

	b := bootstrap.New(config.ServerConfig.Port, config.ServerConfig.LogLevel)
	b.Serve()
}

func listenSignal() {
	if len(collector.NotarySupportedSignals) == 0 {
		return
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, collector.NotarySupportedSignals...)
	select {
	case <-sigs:
		fmt.Println("exitapp,sigs:", sigs)
		os.Exit(0)
	}
}