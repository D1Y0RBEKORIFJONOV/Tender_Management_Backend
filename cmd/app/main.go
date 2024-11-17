package main

import (
	"awesomeProject/internal/app"
	"awesomeProject/internal/config"
	"awesomeProject/logger"
	"fmt"
)

func main() {
	cfg := config.New()
	log := logger.SetupLogger(cfg.LogLevel)
	log.Info("starting app")
	fmt.Println(cfg.RPCPort)
	application := app.NewApp(log, cfg)
	log.Info("starting server")
	forever := make(chan bool)
	go application.HTTPApp.Start()
	<-forever
}
