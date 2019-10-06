package cast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterfaceToSlice(t *testing.T) {
	r := require.New(t)

	testcase := map[string]interface{}{
		"correct":   []interface{}{"a", 3, "x"},
		"malformed": "it's just data",
	}
	_, err := InterfaceToSlice(testcase["correct"])
	r.NoError(err)

	_, err = InterfaceToSlice(testcase["malformed"])
	r.Error(err)

}
