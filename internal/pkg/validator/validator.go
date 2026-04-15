package validator

import (
	v "github.com/go-playground/validator/v10"
)

var validate = v.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// 🔥 helper untuk extract error
func ParseValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if errs, ok := err.(v.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Field() {
			case "Name":
				errors["name"] = "name is required"
			case "Email":
				errors["email"] = "invalid email format"
			case "Password":
				errors["password"] = "password must be at least 6 characters"
			}
		}
	}

	return errors
}
