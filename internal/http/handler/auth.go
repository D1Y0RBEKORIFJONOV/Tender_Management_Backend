package handler

import (
	"awesomeProject/internal/entity"
	authusecase "awesomeProject/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth struct {
	auth authusecase.UserUseCaseImpl
}

func NewAuth(auth authusecase.UserUseCaseImpl) *Auth {
	return &Auth{auth: auth}
}

func (u *Auth) Register(c *gin.Context) {
	var req entity.CreateUsrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := u.auth.RegisterUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
