package middleware

import (
	"awesomeProject/internal/infastructure/token"
	"fmt"
	"github.com/casbin/casbin/v2"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	rateLimiters = map[string]*rate.Limiter{
		"user":         rate.NewLimiter(1, 10),
		"unauthorized": rate.NewLimiter(1, 10),
	}
)

func Middleware(c *gin.Context) {
	allow, err := CheckPermission(c.Request)
	if err != nil {
		if valid, ok := err.(*jwt.ValidationError); ok && valid.Errors == jwt.ValidationErrorExpired {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "token was expired",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	} else if !allow {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
		return
	}

	role, _ := GetRole(c.Request)
	limiter, exists := rateLimiters[role]
	if !exists {
		limiter = rateLimiters["unauthorized"]
	}

	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "rate limit exceeded",
		})
		return
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	id, _ := token.GetIdFromToken(c.Request)
	c.Set("user_id", id)
	email, _ := token.GetEmailFromToken(c.Request)
	c.Set("email", email)
	c.Next()
}

func TimingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start)
	c.Writer.Header().Set("X-Response-Time", duration.String())
}

func CheckPermission(r *http.Request) (bool, error) {
	role, err := GetRole(r)
	if err != nil {
		log.Println("Error while getting role from token: ", err)
		return false, err
	}

	method := r.Method
	path := r.URL.Path

	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		log.Println(err)
	}
	allowed, err := enforcer.Enforce(role, path, method)
	if err != nil {
		log.Println(err)
	}

	fmt.Print(">>>>>>>", allowed)

	return allowed, nil
}

func GetRole(r *http.Request) (string, error) {
	tokenStr := r.Header.Get("Authorization")

	if tokenStr == "" {
		return "unauthorized", nil
	} else if strings.Contains(tokenStr, "Basic") {
		return "unauthorized", nil
	}

	claims, err := token.ExtractClaim(tokenStr)
	if err != nil {
		log.Println("Error while extracting claims: ", err)
		return "unauthorized", err
	}

	return claims["role"].(string), nil
}
