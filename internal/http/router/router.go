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

	router.POST("/register", authhandler.Register)
	router.POST("/login", authhandler.LoginUser)

	router.Use(middleware.Middleware)
	router.Use(middleware.TimingMiddleware)

	tenderGroup := router.Group("")
	{
		tenderGroup.POST("/api/client/tenders", tenderhandler.CreateTender)
		tenderGroup.PUT("/api/client/tenders/:tenderId", tenderhandler.UpdateTenderStatus) 
		tenderGroup.DELETE("/api/client/tenders/:tenderId", tenderhandler.DeleteTender)
		tenderGroup.GET("/api/client/tenders", tenderhandler.GetTenders)
	}

	bidGroup := router.Group("")
	{
		bidGroup.POST("tenders/:id/bids", bidhandler.CreateBid)
		bidGroup.GET("tenders/bids/:id", bidhandler.GetBids)
		bidGroup.PUT("tenders/:id/bids", bidhandler.UpdateBid)
		bidGroup.DELETE("tenders/bids/:id", bidhandler.DeleteBid)
		bidGroup.POST("tenders/:id/award/:bid_id", bidhandler.AnnounceWinner)
	}

	return router
}
