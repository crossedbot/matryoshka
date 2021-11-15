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
		BuildCommands: []string{"make -f %s/Makefile -C %s"},
		RunCommands:   []string{"%s/a.out"},
	},
}

func (c *cLang) Build(mainFile string, otherFiles ...string) (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	c.binDir = filepath.Dir(mainFile)
	build := fmt.Sprintf(c.BuildCommands[0], cwd, c.binDir)
	build = strings.TrimSpace(build)
	parts := strings.Split(build, " ")
	command := parts[0]
	args := parts[1:]
	return execute(CBuildTimeout, command, args...)
}

func (c *cLang) Run(timeout time.Duration) (string, string, error) {
	return execute(timeout, fmt.Sprintf(c.RunCommands[0], c.binDir))
}
