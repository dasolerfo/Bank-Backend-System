package api

import (
	db "simplebank/db/model"

	"github.com/go-playground/validator/v10"
)

var validateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currencyStr, ok := fl.Field().Interface().(string); ok {
		var c db.Currency
		if err := c.Scan(currencyStr); err != nil {
			return false // No és una currency vàlida
		}
		return true // És vàlida
	}
	return false // No és un string, així que no es pot validar
}
