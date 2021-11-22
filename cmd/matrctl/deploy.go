package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"

	"github.com/crossedbot/matryoshka/pkg/deployer/controller"
	"github.com/crossedbot/matryoshka/pkg/runner"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type runCodeCmd struct {
	// Flags
	data              string
	language          string
	os                string
	arch              string
	filepath          string
	content           string
	timeout           int
	preBuildCommands  ArrayFlag
	postBuildCommands ArrayFlag
	preRunCommands    ArrayFlag
	postRunCommands   ArrayFlag
}

func (*runCodeCmd) Name() string {
	return "run-code"
}

func (*runCodeCmd) Synopsis() string {
	return "Run code for a given programming language."
}

func (*runCodeCmd) Usage() string {
	return `run-code [-data <data>] [-language <language>] [-os <os>]
         [-arch <arch>] [-filepath <filepath>] [-directory <dir>]
         [-content <content>] [-timeout <timeout>]
         [-pre-build-command <cmd>]... [-post-build-commands <cmd>]...
	 [-pre-run-command <cmd>]... [-post-run-commands <cmd>]...:
  Run code for a given language.
`
}

func (rc *runCodeCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&rc.data, "data", "", "JSON formatted payload data")
	f.StringVar(&rc.language, "language", "", "Programming language of the content")
	f.StringVar(&rc.os, "os", "debian", "The targeted operating system")
	f.StringVar(&rc.arch, "arch", "amd64", "The targeted architecture")
	f.StringVar(&rc.filepath, "filepath", "", "Location of file that contains the content")
	f.StringVar(&rc.content, "content", "", "Code to run")
	f.IntVar(&rc.timeout, "timeout", 30, "Run timeout in seconds")
	f.Var(&rc.preBuildCommands, "pre-build-command", "Run command before building the payload")
	f.Var(&rc.postBuildCommands, "post-build-command", "Run command after building the payload")
	f.Var(&rc.preRunCommands, "pre-run-command", "Run command before running the payload")
	f.Var(&rc.postRunCommands, "post-run-command", "Run command after running the payload")
}

func (rc *runCodeCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	var payload runner.Payload
	var err error
	// Parse the payload data, or do so through the language, filepath,
	// and/or file content
	if rc.data != "" {
		// Parse the data payload
		payload, err = parseDataPayload(rc.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", rc.Name(), err)
			return subcommands.ExitUsageError
		}
	} else if rc.language != "" && (rc.filepath != "" || rc.content != "") {
		payload.Language = rc.language
		// Parse the filepath for the contents of the payload
		if rc.filepath != "" {
			payloadFile, err := parseFilePayload(rc.filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"%s: %s\n", rc.Name(), err)
				return subcommands.ExitUsageError
			}
			payload.Files = append(payload.Files, payloadFile)
		} else {
			// Else use the provided code content
			payloadFile, err := parseContentPayload(
				payload.Language, rc.content)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"%s: %s\n", rc.Name(), err)
				return subcommands.ExitFailure
			}
			payload.Files = append(payload.Files, payloadFile)
		}
	} else {
		// Exit with error if expectations are not met
		fmt.Fprintf(os.Stderr, "%s: payload data is required\n", rc.Name())
		return subcommands.ExitUsageError
	}

	payload.OperatingSystem = rc.os
	payload.Architecture = rc.arch
	payload.Timeout = rc.timeout
	payload.PreBuildCommands = rc.preBuildCommands
	payload.PostBuildCommands = rc.postBuildCommands
	payload.PreRunCommands = rc.preRunCommands
	payload.PostRunCommands = rc.postRunCommands

	// Create the deployment using the code deployment controller
	result, err := controller.V1().CreateDeployment(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", rc.Name(), err)
		return subcommands.ExitFailure
	}

	// Otherwise we can write the results output to STDOUT
	printOutput(result)

	return subcommands.ExitSuccess
}

func generateRandomString(sz int) (string, error) {
	// TODO this is repeated in the code runner, put this function in a
	// common place
	b := make([]byte, sz)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func parseDataPayload(data string) (runner.Payload, error) {
	var payload runner.Payload
	b := []byte(data)
	if strings.EqualFold(data[:1], "@") {
		// if pointing to a file, retrieve the contents
		contents, err := ioutil.ReadFile(data[1:])
		if err != nil {
			return payload,
				fmt.Errorf("failed to read file '%s'", data[1:])
		}
		b = contents
	}
	// parse and return the payload
	err := json.Unmarshal(b, &payload)
	if err != nil {
		err = fmt.Errorf("failed to parse payload")
	}
	return payload, err
}

func parseFilePayload(f string) (runner.PayloadFile, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return runner.PayloadFile{},
			fmt.Errorf("failed to read file '%s'", f)
	}
	return runner.PayloadFile{
		Name:    filepath.Base(f),
		Content: string(b),
	}, nil
}

func parseContentPayload(lang, content string) (runner.PayloadFile, error) {
	// generate a random name for the file with extension
	filename, err := generateRandomString(runner.KeyLength)
	if err != nil {
		return runner.PayloadFile{},
			fmt.Errorf("failed to generate content filename")
	}
	return runner.PayloadFile{
		Name:    fmt.Sprintf("%s.%s", filename, lang),
		Content: content,
	}, nil
}

func printOutput(result runner.Result) {
	printOutputPlain(result)
}

func printOutputPlain(result runner.Result) {
	for _, v := range result.BuildCommands {
		fmt.Printf("%s+ %s\n", ColorYellow, v.Command)
		if v.Stderr != "" {
			fmt.Fprintf(os.Stderr, "%s%s", ColorRed, v.Stderr)
		}
		if v.Stdout != "" {
			fmt.Fprintf(os.Stdout, "%s%s", ColorReset, v.Stdout)
		}
	}
	for _, v := range result.RunCommands {
		fmt.Printf("%s+ %s\n", ColorYellow, v.Command)
		if v.Stderr != "" {
			fmt.Fprintf(os.Stderr, "%s%s", ColorRed, v.Stderr)
		}
		if v.Stdout != "" {
			fmt.Fprintf(os.Stdout, "%s%s", ColorReset, v.Stdout)
		}
	}
}
