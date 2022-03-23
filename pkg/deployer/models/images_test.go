package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageFilterAdd(t *testing.T) {
	imgFilter := make(ImageFilter)
	key := "HELLO"
	val := "world"
	val2 := "darkness"

	imgFilter.Add(key, val)
	imgFilter.Add(key, val2)
	fValue := imgFilter[strings.ToLower(key)]
	require.Equal(t, 2, len(fValue))
	require.Equal(t, val, fValue[0])
	require.Equal(t, val2, fValue[1])
}

func TestImageFilterSet(t *testing.T) {
	imgFilter := make(ImageFilter)
	key := "HELLO"
	val := "world"
	val2 := "darkness"
	expected := "there"

	imgFilter.Add(key, val)
	imgFilter.Add(key, val2)
	imgFilter.Set(key, expected)
	fValue := imgFilter[strings.ToLower(key)]
	require.Equal(t, 1, len(fValue))
	require.Equal(t, expected, fValue[0])
}

func TestImageFilterGet(t *testing.T) {
	imgFilter := make(ImageFilter)
	key := "HELLO"
	val := "world"
	val2 := "darkness"

	imgFilter.Add(key, val)
	imgFilter.Add(key, val2)
	actual := imgFilter.Get(key)
	require.Equal(t, val, actual)
}

func TestImageFilterDelete(t *testing.T) {
	imgFilter := make(ImageFilter)
	key := "HELLO"
	val := "world"

	imgFilter.Set(key, val)
	imgFilter.Delete(key)
	actual := imgFilter.Get(key)
	require.Empty(t, actual)
}

func TestImageFilterHas(t *testing.T) {
	imgFilter := make(ImageFilter)
	key := "HELLO"
	val := "world"

	imgFilter.Set(key, val)
	actual := imgFilter.Has(key)
	require.True(t, actual)
}

func TestImageSummaryName(t *testing.T) {
	sum := ImageSummary{Repository: "myrepo", Tag: "sometag"}
	expected := fmt.Sprintf("%s:%s", sum.Repository, sum.Tag)
	actual := sum.Name()
	require.Equal(t, expected, actual)
}

func TestImageSummaryLanguage(t *testing.T) {
	sum := ImageSummary{Repository: "repo/c"}
	expected := "c"
	actual := sum.Language()
	require.Equal(t, expected, actual)
}

func TestImageSummaryOperatingSystem(t *testing.T) {
	sum := ImageSummary{Tag: "debian-amd64"}
	expected := "debian"
	actual := sum.OperatingSystem()
	require.Equal(t, expected, actual)
}

func TestImageSummaryArchitecture(t *testing.T) {
	sum := ImageSummary{Tag: "debian-amd64"}
	expected := "amd64"
	actual := sum.Architecture()
	require.Equal(t, expected, actual)
}
