package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/jiny0x01/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
		// check currency is supported
	}
	return false
}
