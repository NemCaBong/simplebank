package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/techschool/simplebank/db/util"
)

// custom validators file in Go

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false

}
