package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GoBuildTimeout = 60 * time.Second
)

// goLang represents the Go programming language
type goLang language

// LanguageGo represents the commands associated with building and running a Go
// program.
var LanguageGo = &goLang{
	// Build commands
	PreBuildCommands:  []string{},
	BuildCommands:     []string{"go build -o a.out ./%s"},
	PostBuildCommands: []string{},
	// Run commands
	PreRunCommands:  []string{},
	RunCommands:     []string{"./a.out"},
	PostRunCommands: []string{},
}

// Build builds a Go program using the given main and additional files. Returns
// the STDOUT and STDERR outputs.
func (g *goLang) Build(mainFile string, otherFiles ...string) (CommandStreams, error) {
	// Set up environment
	commandStreams := CommandStreams{}
	rootDir := filepath.Dir(mainFile)
	rootDir = filepath.Clean(rootDir)
	if err := setGopath(rootDir); err != nil {
		return commandStreams, err
	}
	// Run pre-build commands
	for _, preCmd := range g.PreBuildCommands {
		preCmd = strings.TrimSpace(preCmd)
		command, args := buildCommandLine(preCmd)
		stdout, stderr, err := execute(GoBuildTimeout, command, args...)
		commandStreams = append(
			commandStreams,
			Streams{preCmd, stdout, stderr},
		)
		if err != nil {
			return commandStreams, err
		}
	}
	// Build the binary
	{
		cmd := strings.TrimSpace(fmt.Sprintf(g.BuildCommands[0], rootDir))
		command, args := buildCommandLine(cmd)
		stdout, stderr, err := execute(GoBuildTimeout, command, args...)
		commandStreams = append(commandStreams, Streams{cmd, stdout, stderr})
		if err != nil {
			return commandStreams, err
		}
	}
	// Run post-build
	for _, postCmd := range g.PostBuildCommands {
		postCmd = strings.TrimSpace(postCmd)
		command, args := buildCommandLine(postCmd)
		stdout, stderr, err := execute(GoBuildTimeout, command, args...)
		commandStreams = append(
			commandStreams,
			Streams{postCmd, stdout, stderr},
		)
		if err != nil {
			return commandStreams, err
		}
	}
	return commandStreams, nil
}

// Run the build Go program. Returns the STDOUT and STERR outputs.
func (g *goLang) Run(timeout time.Duration) (CommandStreams, error) {
	commandStreams := CommandStreams{}
	// Run pre-run commands
	for _, preCmd := range g.PreRunCommands {
		preCmd = strings.TrimSpace(preCmd)
		command, args := buildCommandLine(preCmd)
		stdout, stderr, err := execute(GoBuildTimeout, command, args...)
		commandStreams = append(
			commandStreams,
			Streams{preCmd, stdout, stderr},
		)
		if err != nil {
			return commandStreams, err
		}
	}
	// Run the binary
	{
		cmd := strings.TrimSpace(g.RunCommands[0])
		command, args := buildCommandLine(cmd)
		stdout, stderr, err := execute(timeout, command, args...)
		commandStreams = append(commandStreams, Streams{cmd, stdout, stderr})
		if err != nil {
			return commandStreams, err
		}
	}
	// Run post-run commands
	for _, postCmd := range g.PostRunCommands {
		postCmd = strings.TrimSpace(postCmd)
		command, args := buildCommandLine(postCmd)
		stdout, stderr, err := execute(GoBuildTimeout, command, args...)
		commandStreams = append(
			commandStreams,
			Streams{postCmd, stdout, stderr},
		)
		if err != nil {
			return commandStreams, err
		}
	}
	return commandStreams, nil
}

func (g *goLang) SetPreBuildCommands(cmd []string) {
	g.PreBuildCommands = cmd
}

func (g *goLang) SetPostBuildCommands(cmd []string) {
	g.PostBuildCommands = cmd
}

func (g *goLang) SetPreRunCommands(cmd []string) {
	g.PreRunCommands = cmd
}

func (g *goLang) SetPostRunCommands(cmd []string) {
	g.PostRunCommands = cmd
}

func setGopath(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	gopath := filepath.Join(cwd, dir)
	parts := strings.SplitAfterN(gopath, "/src", 2)
	if len(parts) == 2 {
		gopath = filepath.Dir(parts[0])
	}
	os.Setenv("GOPATH", gopath)
	return nil
}
