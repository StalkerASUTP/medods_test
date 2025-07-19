package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("uuid", validateUUID)
}

func validateUUID(fl validator.FieldLevel) bool {
	uuidStr := fl.Field().String()

	if uuidStr == "" {
		return true
	}

	_, err := uuid.Parse(uuidStr)
	return err == nil
}


// Функция для валидации структур
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// Функция для валидации отдельных значений
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}
