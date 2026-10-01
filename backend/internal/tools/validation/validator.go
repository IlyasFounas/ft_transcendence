package validation

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

var (
	upperRegex = regexp.MustCompile(`[A-Z]`)
	lowerRegex = regexp.MustCompile(`[a-z]`)
	numberRegex = regexp.MustCompile(`[0-9]`)
	specialRegex = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func RegisterCustomRules(v *validator.Validate) {
	v.RegisterValidation("strongpwd", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) >= 8 && 
			len(password) <= 50 &&
			upperRegex.MatchString(password) &&
			lowerRegex.MatchString(password) &&
			numberRegex.MatchString(password) &&
			specialRegex.MatchString(password)
	})
}