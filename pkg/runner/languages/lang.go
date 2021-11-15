package languages

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Language represents the interface to programming language.
type Language interface {
	// Build builds the given files. Returning the STDOUT and STDERR output
	// data.
	Build(mainFile string, otherFiles ...string) (string, string, error)

	// Run runs the built program and returns the STDOUT and STDERR output.
	Run(timeout time.Duration) (string, string, error)
}

// language implements the Language interface.
type language struct {
	BuildCommands []string
	RunCommands   []string
}

// Languages is a list of known programming languages.
var Languages = []Language{
	LanguageC,
	LanguageGo,
}

// LanguageStrings is a list of labels for programming languages.
var LanguageStrings = []string{
	"c",
	"go",
}

// ParseLanguage returns the Language for the given label.
func ParseLanguage(v string) (Language, error) {
	for i, s := range LanguageStrings {
		if strings.EqualFold(v, s) && i < len(Languages) {
			return Languages[i], nil
		}
	}
	return nil, fmt.Errorf("unknown language: %s", v)
}

// execute executes the command with the given arguments. Returning the STDOUT
// and STDERR output.
func execute(timeout time.Duration, exe string, args ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(exe, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	processKilled := false
	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
		processKilled = true
	})
	err := cmd.Wait()
	timer.Stop()
	if processKilled {
		err = fmt.Errorf("execution timeout exceeded (%s)",
			timeout.String())
	}
	return stdout.String(), stderr.String(), err
}
