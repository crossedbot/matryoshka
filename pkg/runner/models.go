package runner

import ()

// Payload represents payload's language and files.
type Payload struct {
	Language string        `json:"language"`
	Files    []PayloadFile `json:"files"`

	// Metadata
	OperatingSystem string `json:"operating_system"`
	Architecture    string `json:"architecture"`
	Timeout         int    `json:"timeout"` // in seconds
}

// PayloadFile represents the content and attributes of a file.
type PayloadFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Result represents a the result of returned code.
type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}
