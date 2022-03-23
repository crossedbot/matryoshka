package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/crossedbot/matryoshka/pkg/deployer/models"

	"github.com/stretchr/testify/require"
)

func TestFormatSize(t *testing.T) {
	sz := int64(100)
	expected := "100 B"
	actual := formatSize(sz)
	require.Equal(t, expected, actual)

	sz = int64(1000)
	expected = "1.0 kB"
	actual = formatSize(sz)
	require.Equal(t, expected, actual)

	sz = int64(1000000)
	expected = "1.0 MB"
	actual = formatSize(sz)
	require.Equal(t, expected, actual)

	sz = int64(1000000000)
	expected = "1.0 GB"
	actual = formatSize(sz)
	require.Equal(t, expected, actual)
}

func TestPrintImagesPlain(t *testing.T) {
	created := "2011-03-02T22:11:00Z"
	createdTime, err := time.Parse(time.RFC3339, created)
	require.Nil(t, err)
	imageSummaries := []models.ImageSummary{{
		ID:         "1a2b3c4d5e6f",
		Repository: "matryoshka/test",
		Tag:        "debian-amd64",
		CreatedAt:  createdTime,
		Size:       int64(630600000),
	}}
	expected := fmt.Sprintf(
		"ID=%s\nRepository=%s\nTag=%s\nCreatedAt=%s\nSize=%s\n\n",
		imageSummaries[0].ID, imageSummaries[0].Repository,
		imageSummaries[0].Tag,
		imageSummaries[0].CreatedAt.Format(time.RFC3339),
		formatSize(imageSummaries[0].Size),
	)

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printImagesPlain(imageSummaries)

	w.Close()
	actual, err := ioutil.ReadAll(r)
	require.Nil(t, err)
	os.Stdout = stdout
	require.Equal(t, expected, string(actual))
}

func TestPrintImagesJson(t *testing.T) {
	created := "2011-03-02T22:11:00Z"
	createdTime, err := time.Parse(time.RFC3339, created)
	require.Nil(t, err)
	imageSummaries := []models.ImageSummary{{
		ID:         "1a2b3c4d5e6f",
		Repository: "matryoshka/test",
		Tag:        "debian-amd64",
		CreatedAt:  createdTime,
		Size:       int64(630600000),
	}}

	expected, err := json.MarshalIndent(imageSummaries, "", "    ")
	require.Nil(t, err)

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printImagesJson(imageSummaries)

	w.Close()
	actual, err := ioutil.ReadAll(r)
	require.Nil(t, err)
	os.Stdout = stdout
	require.Equal(t, string(expected), string(actual))
}
