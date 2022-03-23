package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crossedbot/matryoshka/pkg/runner"
)

func TestParseDataPayload(t *testing.T) {
	lang := "c"
	files := []runner.PayloadFile{{
		Name:    "main.c",
		Content: "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n",
	}}
	bFiles, err := json.Marshal(files)
	require.Nil(t, err)
	os := "debian"
	arch := "amd64"
	to := 30
	expected := runner.Payload{
		Language:        lang,
		Files:           files,
		OperatingSystem: os,
		Architecture:    arch,
		Timeout:         to,
	}
	data := fmt.Sprintf(`{
	"language": "%s",
	"files": %s,
	"operating_system": "%s",
	"architecture": "%s",
	"timeout": %d
}`, lang, bFiles, os, arch, to)

	actual, err := parseDataPayload(data)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestParseFilePayload(t *testing.T) {
	name := "test.c"
	content := "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n"
	expected := runner.PayloadFile{
		Name:    name,
		Content: content,
	}

	actual, err := parseFilePayload("testdata/c/test.c")
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestVisitDirectory(t *testing.T) {
	d := "testdata/c"
	files := runner.PayloadFiles{}
	file := "testdata/c/test.c"
	finfo, err := os.Lstat(file)
	require.Nil(t, err)
	expected := runner.PayloadFile{
		Name:    "test.c",
		Path:    "c/",
		Content: "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n",
	}

	walkFn := visitDirectory(d, &files)
	err = walkFn(file, finfo, nil)
	require.Nil(t, err)
	require.Equal(t, 1, len(files))
	require.Equal(t, expected, files[0])
}

func TestParseDirectoryPayload(t *testing.T) {
	d := "testdata/c"
	expected := []runner.PayloadFile{{
		Name:    "test.c",
		Path:    "c/",
		Content: "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n",
	}}
	actual, err := parseDirectoryPayload(d)
	require.Nil(t, err)
	require.Equal(t, 1, len(actual))
	require.Equal(t, expected, actual)
}

func TestParseContentPayload(t *testing.T) {
	lang := "c"
	content := "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n"

	actual, err := parseContentPayload(lang, content)
	require.Nil(t, err)
	require.Equal(t, "."+lang, filepath.Ext(actual.Name))
	require.Equal(t, content, actual.Content)
}
