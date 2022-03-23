package runner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayloadFilesLen(t *testing.T) {
	files := PayloadFiles{{Name: "test.c"}}
	expected := 1
	actual := files.Len()
	require.Equal(t, expected, actual)
}

func TestPayloadFilesLess(t *testing.T) {
	files := PayloadFiles{
		{Name: "a.c", Path: "src"},
		{Name: "b.c", Path: "src"},
	}

	// By name
	actual := files.Less(0, 1)
	require.True(t, actual)
	actual = files.Less(1, 0)
	require.False(t, actual)

	// By depth
	files[0].Path = "src/lib"
	actual = files.Less(0, 1)
	require.False(t, actual)
	actual = files.Less(1, 0)
	require.True(t, actual)
}

func TestPayloadFilesSwap(t *testing.T) {
	name1 := "one"
	name2 := "two"
	files := PayloadFiles{{Name: name1}, {Name: name2}}
	files.Swap(0, 1)
	require.Equal(t, name2, files[0].Name)
	require.Equal(t, name1, files[1].Name)
}
