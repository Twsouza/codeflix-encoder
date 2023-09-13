package utils_test

import (
	"testing"

	"github.com/Twsouza/codeflix-encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJsonValid(t *testing.T) {
	json := `{"name":"John","age":30,"car":null}`
	err := utils.IsJson(json)
	require.Nil(t, err)
}

func TestIsJsonInvalid(t *testing.T) {
	json := `test`
	err := utils.IsJson(json)
	require.Error(t, err)
}
