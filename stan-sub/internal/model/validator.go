package model

import (
	"github.com/go-playground/validator"
)

func ValidateOrder(order Order) error {
	validate := validator.New()
	return validate.Struct(order)
}
