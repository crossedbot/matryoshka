package runner

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/crossedbot/common/golang/logger"

	"github.com/crossedbot/matryoshka/pkg/runner/languages"
)

const (
	KeyLength = 32

	// Defaults
	DefaultRunTimeout = 30
)

// Runner represents a code runner.
type Runner interface {
	// Run starts the process thread for program payloads. If indicated, the
	// process thread returns after running once
	Run(once bool) error

	// Stop stops the running processing thread of the runner.
	Stop() error
}

// runner implements the Runner interface.
type runner struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// New returns a new code runner.
func New(ctx context.Context) Runner {
	ctx, cancel := context.WithCancel(ctx)
	return &runner{ctx, cancel}
}

func (r *runner) Run(once bool) error {
	stdin := make(chan []byte)
	defer close(stdin)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			stdin <- scanner.Bytes()
		}
	}()
	r.process(stdin, once)
	return nil
}

func (r *runner) Stop() error {
	r.cancel()
	return nil
}

// process processes incoming input from STDIN, returning when context is done
// or running once when set.
func (r *runner) process(stdin chan []byte, once bool) {
	run := true
	for run {
		select {
		case <-r.ctx.Done():
			run = false
		case in := <-stdin:
			result := runCode(r.ctx, in)
			if result.Error != "" {
				logger.Error("Error running code:", result.Error)
			}
			b, err := json.Marshal(result)
			if err != nil {
				logger.Error("Error marshalling result:", err.Error())
			}
			// Print result to STDOUT for reading from listeners
			// (E.g. deployer)
			fmt.Printf(string(b))
			run = !once
		}
	}
}

// runCode writes the payload, builds the code, and runs it. Returning the
// result when the code exits.
func runCode(ctx context.Context, data []byte) Result {
	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return Result{Error: err.Error()}
	}
	// write the payload to files and parse payload language
	paths, err := writeFiles(payload.Files)
	if err != nil {
		return Result{Error: err.Error()}
	}
	lang, err := languages.ParseLanguage(payload.Language)
	if err != nil {
		return Result{Error: err.Error()}
	}
	if payload.Timeout < 1 {
		payload.Timeout = DefaultRunTimeout
	}
	timeout := time.Duration(payload.Timeout) * time.Second
	lang.SetPreBuildCommands(payload.PreBuildCommands)
	lang.SetPostBuildCommands(payload.PostBuildCommands)
	lang.SetPreRunCommands(payload.PreRunCommands)
	lang.SetPostRunCommands(payload.PostRunCommands)
	// build the code using the written files for the selected language
	buildStreams, err := lang.Build(paths[0], paths[1:]...)
	result := Result{
		BuildCommands: buildStreams,
		RunCommands:   languages.CommandStreams{},
		Error:         "",
	}
	if err != nil {
		result.Error = fmt.Sprintf(
			"Error while building: %s",
			err.Error(),
		)
		return result
	}
	// run the built code and return the result
	runStreams, err := lang.Run(timeout)
	result.RunCommands = runStreams
	if err != nil {
		result.Error = fmt.Sprintf(
			"Error while running: %s",
			err.Error(),
		)
	}
	return result
}

// writeFiles writes the given files to a created temporary subdirectory of the
// current working directory. Returns the path for each file written.
func writeFiles(files []PayloadFile) ([]string, error) {
	// Get current working directory and create a subdirectory
	cwd, err := os.Getwd()
	if err != nil {
		return []string{}, err
	}
	dir, err := createSubdirectory(cwd)
	if err != nil {
		return []string{}, err
	}
	// For each payload file, write it to the subdirectory we just created
	var paths []string
	for _, file := range files {
		// Create parent directory before writing the file itself
		location := filepath.Join(dir, file.Path)
		if err := os.MkdirAll(location, 0700); err != nil {
			return []string{}, err
		}
		// Write the file to the directory
		location = filepath.Join(location, file.Name)
		err := ioutil.WriteFile(location, []byte(file.Content), 0644)
		if err != nil {
			return []string{}, err
		}
		// Track the relative path to the file that was just written
		rp, err := filepath.Rel(cwd, location)
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, rp)
	}
	// Return the paths to all files written
	return paths, nil
}

// createSubdirectory creates a randomly generated subdirectory of the 'tmp' and
// the given parent directory.
func createSubdirectory(parent string) (string, error) {
	// Create a randomly named subdirectory, since we assume it to be
	// temporary
	subdir, err := generateRandomString(KeyLength)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(parent, "tmp", subdir, "src")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

// generateRandomString generates a random string of size 'sz'.
func generateRandomString(sz int) (string, error) {
	b := make([]byte, sz)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
