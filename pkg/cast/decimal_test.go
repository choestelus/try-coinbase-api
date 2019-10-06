package cast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterfaceStringToDecimal(t *testing.T) {
	r := require.New(t)
	correctInput, err := InterfaceStringToDecimal("1410.2516")
	r.NoError(err)
	r.Equal(correctInput.String(), "1410.2516")

	_, err = InterfaceStringToDecimal("justice")
	r.Error(err)
}
