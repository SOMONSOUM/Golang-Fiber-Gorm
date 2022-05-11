package util

import (
	"regexp"

	"github.com/SOMONSOUM/go-fiber/models"
	valid "github.com/asaskevich/govalidator"
)

func IsEmpty(str string) (bool, string) {
	if valid.HasWhitespaceOnly(str) && str != "" {
		return true, "Must not be empty"
	}

	return false, ""
}

func ValidateRegister(user *models.User) *models.UserErrors {
	e := &models.UserErrors{}
	e.Err, e.Username = IsEmpty(user.Username)

	if !valid.IsEmail(user.Email) {
		e.Err, e.Email = true, "Must be a valid email"
	}

	re := regexp.MustCompile("\\d") // regex check for at least one integer in string
	if !(len(user.Password) >= 8 && valid.HasLowerCase(user.Password) && valid.HasUpperCase(user.Password) && re.MatchString(user.Password)) {
		e.Err, e.Password = true, "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"
	}

	return e
}
