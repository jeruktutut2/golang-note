package modelrequest

type LoginRequest struct {
	// Username string `json:"username" validate:"required,usernamevalidator"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwordvalidator"`
}
