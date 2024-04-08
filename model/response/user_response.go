package modelresponse

import modelentity "golang-note/model/entity"

type LoginResponse struct {
	Id       uint32 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Utc      string `json:"utc"`
}

func ToLoginResponse(user modelentity.User) (loginResponse LoginResponse) {
	// var loginResponse LoginResponse
	loginResponse.Id = uint32(user.Id.Int32)
	loginResponse.Username = user.Username.String
	loginResponse.Email = user.Email.String
	loginResponse.Utc = user.Utc.String
	return
}
