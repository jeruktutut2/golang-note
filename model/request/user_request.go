package modelrequest

type LoginRequest struct {
	Username string `json:"username" validate:"required,usernamevalidator"`
	Password string `json:"password" validate:"required,passwordvalidator"`
}
