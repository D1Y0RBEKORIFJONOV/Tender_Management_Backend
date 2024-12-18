package main

import (
	"awesomeProject/internal/app"
	"awesomeProject/internal/config"
	"awesomeProject/internal/infastructure/repository/mongodb"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	notifactionservice "awesomeProject/internal/service/notifaction"
	notificationusecase "awesomeProject/internal/usecase/notification"
	"awesomeProject/logger"
	"awesomeProject/websocket"
	"fmt"
	"log"
	"net/http"
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
	go openwebsockets()
	<-forever

}

func openwebsockets() {
	cfg := config.New()
	logger := logger.SetupLogger(cfg.LogLevel)
	db, err := mongodb.NewMongoDB(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
	}
	redisC := redisrepository.NewRedis(*cfg)
	dbUsecase := notificationusecase.NewNotificationRepository(db)

	notification := notifactionservice.NewNotification(logger, dbUsecase, redisC)
	notificationService := notificationusecase.NewNotificationUseCase(notification)
	server := websocket.NewServer(notificationService)

	http.HandleFunc("/ws", server.HandlerNotification)
	log.Fatal(http.ListenAndServe(":9005", nil))
}
