package cast

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/shopspring/decimal"
)

// InterfaceStringToDecimal convert interface{} directly into decimal type
// if supplied interface is string
func InterfaceStringToDecimal(i interface{}) (decimal.Decimal, error) {
	str, ok := i.(string)
	if !ok {
		return decimal.Zero, errors.Wrap(fmt.Errorf("type assertion failed"), "interface -> string")
	}
	return decimal.NewFromString(str)
}
