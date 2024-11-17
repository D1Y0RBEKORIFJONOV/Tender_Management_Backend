package httpapp

import (
	"awesomeProject/internal/http/router"
	authusecase "awesomeProject/internal/usecase/auth"
	bidusecase "awesomeProject/internal/usecase/bid"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	Port   string
	Server *gin.Engine
}

func NewApp(logger *slog.Logger,
	port string,
	auth *authusecase.UserUseCaseImpl,
	bid *bidusecase.BidUseCaseIml,
	tender *tenderusecase.TenderUseCaseIml) *App {

	server := router.Router(*auth, *tender, *bid)
	return &App{
		Logger: logger,
		Port:   port,
		Server: server,
	}
}

func (app *App) Start() {
	const op = "app.Start"
	log := app.Logger.With(
		slog.String(op, "Starting server"),
		slog.String("port", app.Port))
	log.Info("Starting server")
	err := app.Server.SetTrustedProxies(nil)
	if err != nil {
		log.Error("Error setting trusted proxies", "error", err)
		return
	}
	err = app.Server.Run(app.Port)
	if err != nil {
		log.Error("Error starting server", "error", err)
		return
	}
}
