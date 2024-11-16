package entity

type (
	User struct {
		ID       string `json:"id" bson:"_id,omitempty"`
		Username string `json:"username" bson:"username"`
		Password string `json:"password" bson:"password"`
		Email    string `json:"email" bson:"email"`
		Role     string `json:"role" bson:"role"`
	}
	CreateUsrRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
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
)
