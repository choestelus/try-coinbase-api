package cast

import (
	"fmt"

	"emperror.dev/errors"
)

func InterfaceToSlice(i interface{}) ([]interface{}, error) {
	si, ok := i.([]interface{})
	if !ok {
		return nil, errors.Wrap(fmt.Errorf("type assertion failed"), "interface->[]interface")
	}
	return si, nil
}
