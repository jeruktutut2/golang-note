package setup

import (
	"golang-note/helper"

	"github.com/go-playground/validator/v10"
)

func Validate() (validate *validator.Validate) {
	validate = validator.New()
	helper.UsernameValidator(validate)
	helper.PasswordValidator(validate)
	helper.TelephoneValidator(validate)
	return
}
