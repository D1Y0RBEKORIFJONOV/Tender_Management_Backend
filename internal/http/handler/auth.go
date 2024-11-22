package handler

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/entity"
	authusecase "awesomeProject/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

type Auth struct {
	auth authusecase.UserUseCaseImpl
}

func NewAuth(auth authusecase.UserUseCaseImpl) *Auth {
	return &Auth{auth: auth}
}

// @title Artisan Connect
// @version 1.0
// @description This is a sample server for a restaurant reservation system.
// @host 52.59.220.158:9006
// @BasePath        /
// @schemes         https
// @securityDefinitions.apiKey ApiKeyAuth
// @in              header
// @name            Authorization

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
// @Success 200 {object} entity.User
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /register [post]
func (u *Auth) Register(c *gin.Context) {
	var req entity.CreateUsrRequest

	// Bind the incoming JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle missing fields or incorrect JSON
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email cannot be empty",
		})
		return
	}

	// Check if the email or username is empty
	if req.Email == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email cannot be empty",
		})
		return
	}

	// Register the user
	token, err := u.auth.RegisterUser(c.Request.Context(), req)
	if err != nil {
		// Handle specific error cases from RegisterUser function
		switch {
		case strings.Contains(err.Error(), "Email already exists"):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email already exists",
			})
		case strings.Contains(err.Error(), "invalid email format"):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid email format",
			})
		case strings.Contains(err.Error(), "invalid role"):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid role",
			})
		default:
			// For all other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Registration failed",
			})
		}
		return
	}

	// If registration is successful, return the token
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Simple regex to validate email format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// LoginUser godoc
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
