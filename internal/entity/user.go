package entity

import "time"

type (
	User struct {
		ID       string `json:"id" bson:"_id,omitempty"`
		Username string `json:"username" bson:"username"`
		Password string `json:"password" bson:"password"`
		Email    string `json:"email" bson:"email"`
		Role     string `json:"role" bson:"role"`
		Token    string `json:"token" bson:"token"`
	}
	CreateUsrRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	SaveRegisRequest struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		Role       string `json:"role"`
		SecretCode string `json:"secret_code"`
	}
	VerifyUserRequest struct {
		Email      string `json:"email" redis:"email"`
		SecretCode string `json:"secret_code" redis:"secret_code"`
	}
	VerifyUserResponse struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Token string `json:"token"`
	}

	EmailNotificationReq struct {
		SenderName string    `json:"sender_name" bson:"sender_name"`
		SenderAt   time.Time `json:"sender_at" bson:"sender_at"`
		Tittle     string    `json:"Tittle" bson:"Tittle"`
		Content    string    `json:"content" bson:"content"`
		Recipient  []string  `json:"recipient" bson:"recipient"`
	}
)
