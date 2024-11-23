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

	if err := c.ShouldBindJSON(&req); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email cannot be empty",
		})
		return
	}
	if !(req.Role == "client" || req.Role == "contractor") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid role"})
		return
	}

	if req.Email == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email cannot be empty",
		})
		return
	}
	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email format"})
		return
	}
	token, err := u.auth.RegisterUser(c.Request.Context(), req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Email already exists"):
			c.JSON(400, gin.H{
				"message": "Email already exists",
			})
		case strings.Contains(err.Error(), "invalid role"):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid role",
			})
		default:
			c.JSON(500, gin.H{
				"message": "Registration failed",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func isValidEmail(email string) bool {
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
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username and password are required"})
		return
	}

	message, err := u.auth.LoginUser(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "Invalid username or password":
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password",
			})
			return
		case "User not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
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
