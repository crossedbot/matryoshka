package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	CBuildTimeout = 60 * time.Second
)

type cLang struct {
	language
	binDir string
}

var LanguageC = &cLang{
	language: language{
		// Build commands
		PreBuildCommands:  []string{},
		BuildCommands:     []string{"make -f %s/Makefile -C %s"},
		PostBuildCommands: []string{},
		// Run commands
		PreRunCommands:  []string{},
		RunCommands:     []string{"%s/a.out"},
		PostRunCommands: []string{},
	},
}

func (c *cLang) Build(mainFile string, otherFiles ...string) (CommandStreams, error) {
	commandStreams := CommandStreams{}
	// Run pre-build commands
	for _, preCmd := range c.PreBuildCommands {
		preCmd = strings.TrimSpace(preCmd)
		command, args := buildCommandLine(preCmd)
		stdout, stderr, err := execute(CBuildTimeout, command, args...)
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
		cwd, err := os.Getwd()
		if err != nil {
			return commandStreams, err
		}
		c.binDir = filepath.Dir(mainFile)
		cmd := fmt.Sprintf(c.BuildCommands[0], cwd, c.binDir)
		cmd = strings.TrimSpace(cmd)
		command, args := buildCommandLine(cmd)
		stdout, stderr, err := execute(CBuildTimeout, command, args...)
		commandStreams = append(commandStreams, Streams{cmd, stdout, stderr})
		if err != nil {
			return commandStreams, err
		}
	}
	// Run post-build commands
	for _, postCmd := range c.PostBuildCommands {
		postCmd = strings.TrimSpace(postCmd)
		command, args := buildCommandLine(postCmd)
		stdout, stderr, err := execute(CBuildTimeout, command, args...)
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

func (c *cLang) Run(timeout time.Duration) (CommandStreams, error) {
	commandStreams := CommandStreams{}
	// Run pre-run commands
	for _, preCmd := range c.PreRunCommands {
		preCmd = strings.TrimSpace(preCmd)
		command, args := buildCommandLine(preCmd)
		stdout, stderr, err := execute(CBuildTimeout, command, args...)
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
		cmd := fmt.Sprintf(c.RunCommands[0], c.binDir)
		command, args := buildCommandLine(cmd)
		stdout, stderr, err := execute(timeout, command, args...)
		commandStreams = append(
			commandStreams,
			Streams{cmd, stdout, stderr},
		)
		if err != nil {
			return commandStreams, err
		}
	}
	// Run post-run commands
	for _, postCmd := range c.PostRunCommands {
		postCmd = strings.TrimSpace(postCmd)
		command, args := buildCommandLine(postCmd)
		stdout, stderr, err := execute(CBuildTimeout, command, args...)
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

func (c *cLang) SetPreBuildCommands(cmd []string) {
	c.PreBuildCommands = cmd
}

func (c *cLang) SetPostBuildCommands(cmd []string) {
	c.PostBuildCommands = cmd
}

func (c *cLang) SetPreRunCommands(cmd []string) {
	c.PreRunCommands = cmd
}

func (c *cLang) SetPostRunCommands(cmd []string) {
	c.PostRunCommands = cmd
}
