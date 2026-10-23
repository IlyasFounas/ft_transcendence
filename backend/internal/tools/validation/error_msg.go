package validation

import (
	"github.com/go-playground/validator/v10"
)

func FormatErrors(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return "Invalid input."
	}

	fieldError := validationErrors[0]

	switch fieldError.Field() {
	case "Email":
		switch fieldError.Tag() {
		case "required":
			return "Email is required."
		case "max":
			return "Email must be less than or equal to 254 characters."
		case "email":
			return "Please enter a valid email address."
		}

	case "Login":
		switch fieldError.Tag() {
		case "required":
			return "Login ID is required."
		case "min":
			return "Login ID must be at least 3 characters."
		case "max":
			return "Login ID must be at most 30 characters."
		case "acceptChar":
			return "Login ID can only contain letters, numbers, underscores, and hyphens."
		}

	case "Username":
		switch fieldError.Tag() {
		case "required":
			return "Username is required."
		case "min":
			return "Username must be at least 3 characters."
		case "max":
			return "Username must be at most 30 characters."
		case "acceptChar":
			return "Username can only contain letters, numbers, underscores, and hyphens."
		}

	case "Password":
		switch fieldError.Tag() {
		case "required":
			return "Password is required."
		case "strongpwd":
			return "Password must contain at least 8 to 50 characters, one uppercase letter, one lowercase letter, one number, and one special character."
		}
	}

	return "Invalid input."
}
