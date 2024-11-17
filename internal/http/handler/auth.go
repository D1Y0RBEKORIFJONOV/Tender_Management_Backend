package handler

import (
	"awesomeProject/internal/entity"
	authusecase "awesomeProject/internal/usecase/auth"
	"net/http"
	_ "awesomeProject/docs"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	auth authusecase.UserUseCaseImpl
}

func NewAuth(auth authusecase.UserUseCaseImpl) *Auth {
	return &Auth{auth: auth}
}

// @Tender
// @version 1.0
// @description This is a sample server for a Tender  system.
// @securityDefinitions.apikey Bearer
// @in 				header
// @name Authorization
// @description Enter the token in the format `Bearer {token}`
// @host localhost:9006
// @BasePath /

// Register godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body entity.CreateUsrRequest true "User request body"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /register [post]
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

// VerifyCode godoc
// @Summary Verify a user code
// @Description Verify the user code sent to the user's email
// @Tags user
// @Accept json
// @Produce json
// @Param verify body entity.VerifyUserRequest true "Verification request body"
// @Success 200 {object} entity.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /verify-code [post]
func (u *Auth) VerifyUser(c *gin.Context) {
	var req entity.VerifyUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := u.auth.VerifyUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// Login godoc
// @Summary User login
// @Description Log in a user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param login body entity.LoginRequest true "Login request body"
// @Success 200 {object} entity.LoginResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func (u *Auth) LoginUser(c *gin.Context) {
	var req  entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := u.auth.LoginUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : message})
}


