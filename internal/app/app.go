package app

import (
	httpapp "awesomeProject/internal/app/http"
	"awesomeProject/internal/config"
	"awesomeProject/internal/infastructure/repository/mongodb"
	"awesomeProject/internal/infastructure/repository/postgres/bids"
	"awesomeProject/internal/infastructure/repository/postgres/tenders"
	userrepo "awesomeProject/internal/infastructure/repository/postgres/user"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	"awesomeProject/internal/service/auth"
	"awesomeProject/internal/service/bid"
	notifactionservice "awesomeProject/internal/service/notifaction"
	"awesomeProject/internal/service/tender"
	authusecase "awesomeProject/internal/usecase/auth"
	bidusecase "awesomeProject/internal/usecase/bid"
	notificationusecase "awesomeProject/internal/usecase/notification"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"log"
	"log/slog"
)

type App struct {
	HTTPApp *httpapp.App
}

func NewApp(logger *slog.Logger, cfg *config.Config) *App {
	redisDb := redisrepository.NewRedis(*cfg)

	dbAuth, err := userrepo.NewUserRepository()
	if err != nil {
		log.Fatalln(err)
	}
	mongoDb, err := mongodb.NewMongoDB(cfg, logger)
	if err != nil {
		panic(err)
	}

	notificationDb := notificationusecase.NewNotificationRepository(mongoDb)
	notification := notifactionservice.NewNotification(logger, notificationDb, redisDb)
	notificationUseCase := notificationusecase.NewNotificationUseCase(notification)

	authServiceDb := authusecase.NewAuthDbUseCase(dbAuth)
	authService := auth.NewAuth(logger, redisDb, authServiceDb, notificationUseCase)
	authUseCase := authusecase.NewUserUseCase(authService)

	dbTender := tenders.NewTenderRepository()
	tenderDbUseCase := tenderusecase.NewTenderRepository(dbTender)
	tenderService := tender.NewTender(logger, *tenderDbUseCase)
	tenderUseCase := tenderusecase.NewTenderUseCase(tenderService)

	dbBid := bids.NewBidRepository()
	bidsDbUseCase := bidusecase.NewBidUseCaseIml(dbBid)
	bidsService := bid.NewBid(logger, bidsDbUseCase, tenderDbUseCase, notificationUseCase)
	bidUseCase := bidusecase.NewBidUseCaseIml(bidsService)

	server := httpapp.NewApp(logger, cfg.RPCPort, authUseCase, bidUseCase, tenderUseCase)

	return &App{
		HTTPApp: server,
	}
}
