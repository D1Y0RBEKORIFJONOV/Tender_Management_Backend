package router

import (
	"awesomeProject/internal/http/handler"
	"awesomeProject/internal/http/middleware"
	authusecase "awesomeProject/internal/usecase/auth"
	bidusecase "awesomeProject/internal/usecase/bid"
	tenderusecase "awesomeProject/internal/usecase/tender"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(auth authusecase.UserUseCaseImpl, tender tenderusecase.TenderUseCaseIml, bid bidusecase.BidUseCaseIml) *gin.Engine {
	authhandler := handler.NewAuth(auth)
	tenderhandler := handler.NewTender(tender)
	bidhandler := handler.NewBid(bid)

	router := gin.Default()
	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Use(middleware.Middleware)
	router.Use(middleware.TimingMiddleware)

	authGroup := router.Group("/")
	{
		authGroup.POST("register", authhandler.Register)
		authGroup.POST("login", authhandler.LoginUser)
	}

	tenderGroup := router.Group("/")
	{
		tenderGroup.POST("tenders", tenderhandler.CreateTender)
		tenderGroup.GET("tenders", tenderhandler.GetTenders)
		tenderGroup.PUT("tenders/:id", tenderhandler.UpdateTenderStatus)
		tenderGroup.DELETE("tenders/:id", tenderhandler.DeleteTender)
	}

	bidGroup := router.Group("/")
	{
		bidGroup.POST("tenders/bids:id", bidhandler.CreateBid)
		bidGroup.GET("tenders/:id/bids", bidhandler.GetBids)
		bidGroup.PUT("tenders/:id/bids/:id", bidhandler.UpdateBid)
		bidGroup.DELETE("tenders/:id/bids/:id", bidhandler.DeleteBid)
		bidGroup.POST("tenders/:id/bids/winner", bidhandler.AnnounceWinner)
	}

	return router
}
