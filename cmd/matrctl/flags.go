package main

import (
	"fmt"
	"strings"
)

type ArrayFlag []string

func (a *ArrayFlag) String() string {
	return strings.Join(*a, " ")
}

func (a *ArrayFlag) Set(v string) error {
	*a = append(*a, v)
	return nil
}

type OutputFormat string

const (
	OutputFormatPlain = "plain"
	OutputFormatJson  = "json"
)

var OutputFormats = []OutputFormat{
	OutputFormatPlain,
	OutputFormatJson,
}

func (o *OutputFormat) String() string {
	return string(*o)
}

func (o *OutputFormat) Set(v string) error {
	var err error
	*o, err = GetOutputFormat(v)
	if err != nil {
		return err
	}
	return nil
}

func GetOutputFormat(s string) (OutputFormat, error) {
	for _, v := range OutputFormats {
		if strings.EqualFold(s, v.String()) {
			return v, nil
		}
	}
	return "", fmt.Errorf("unkown output format \"%s\"", s)
}
