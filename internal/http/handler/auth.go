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

// Register godoc
// @Summary Create a new user
// @Description Register a new user with email, password, role, and username
// @Tags user
// @Accept json
// @Produce json
// @Param user body entity.CreateUsrRequest true "User registration request body"
// @Success 201 {object}
// @Failure 400 {object}
// @Failure 409 {object}
// @Failure 500 {object}
// @Router /register [post]
func (u *Auth) Register(c *gin.Context) {
	var req entity.CreateUsrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email cannot be empty",
		})
		return
	}

	user, err := u.auth.RegisterUser(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": user.Token,
	})
}

// LoginUser godoc
// @Summary Log in a user
// @Description Log in a user with email/username and password
// @Tags user
// @Accept json
// @Produce json
// @Param login body entity.LoginRequest true "Login request body"
// @Success 200 {object} map[string]interface{} "Successful login response containing a token"
// @Failure 400 {object} map[string]interface{} "Validation error or missing fields"
// @Failure 401 {object} map[string]interface{} "Invalid username or password"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /login [post]
func (u *Auth) LoginUser(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password are required",
		})
		return
	}

	message, err := u.auth.LoginUser(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password",
			})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": message,
	})
}
