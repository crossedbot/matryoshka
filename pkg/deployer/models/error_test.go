package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	err := Error{
		Code:    123,
		Message: "hello world",
	}
	expected := fmt.Sprintf("%d: %s", err.Code, err.Message)
	actual := err.Error()
	require.Equal(t, expected, actual)
}
