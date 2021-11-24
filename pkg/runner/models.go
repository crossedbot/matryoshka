package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/crossedbot/matryoshka/pkg/runner/languages"
)

// Payload represents payload's language and files.
type Payload struct {
	Language string        `json:"language"`
	Files    []PayloadFile `json:"files"`

	// Metadata
	OperatingSystem string `json:"operating_system"`
	Architecture    string `json:"architecture"`
	Timeout         int    `json:"timeout"` // in seconds

	// Commands
	PreBuildCommands  []string `json:"pre_build_commands"`
	PostBuildCommands []string `json:"post_build_commands"`
	PreRunCommands    []string `json:"pre_run_commands"`
	PostRunCommands   []string `json:"post_run_commands"`
}

// Result represents a the result of returned code.
type Result struct {
	BuildCommands languages.CommandStreams `json:"build_commands"`
	RunCommands   languages.CommandStreams `json:"run_commands"`
	Error         string                   `json:"error"`
}

// PayloadFile represents the content and attributes of a file.
type PayloadFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type PayloadFiles []PayloadFile

func (pf PayloadFiles) Len() int { return len(pf) }

func (pf PayloadFiles) Less(i, j int) bool {
	fpath1 := filepath.Join(pf[i].Path, pf[i].Name)
	depth1 := len(strings.Split(fpath1, string(os.PathSeparator)))
	fpath2 := filepath.Join(pf[j].Path, pf[j].Name)
	depth2 := len(strings.Split(fpath2, string(os.PathSeparator)))
	return depth1 < depth2
}

func (pf PayloadFiles) Swap(i, j int) { pf[i], pf[j] = pf[j], pf[i] }
