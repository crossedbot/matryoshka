package languages

import (
	"fmt"
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
	BuildCommands: []string{"go build -o a.out %s"},
	RunCommands:   []string{"./a.out"},
}

// Build builds a Go program using the given main and additional files. Returns
// the STDOUT and STDERR outputs.
func (g *goLang) Build(mainFile string, otherFiles ...string) (string, string, error) {
	build := fmt.Sprintf(g.BuildCommands[0], mainFile)
	build = strings.TrimSpace(build)
	parts := strings.Split(build, " ")
	command := parts[0]
	args := append(parts[1:], otherFiles...)
	return execute(GoBuildTimeout, command, args...)
}

// Run the build Go program. Returns the STDOUT and STERR outputs.
func (g *goLang) Run(timeout time.Duration) (string, string, error) {
	return execute(timeout, g.RunCommands[0])
}
