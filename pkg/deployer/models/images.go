package models

import (
	"fmt"
	"strings"
	"time"
)

type ImageFilter map[string][]string

func (filter ImageFilter) Add(key, val string) {
	key = strings.ToLower(key)
	filter[key] = append(filter[key], val)
}

func (filter ImageFilter) Set(key, val string) {
	key = strings.ToLower(key)
	filter[key] = []string{val}
}

func (filter ImageFilter) Get(key string) string {
	key = strings.ToLower(key)
	if filter == nil {
		return ""
	}
	v := filter[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (filter ImageFilter) Delete(key string) {
	key = strings.ToLower(key)
	delete(filter, key)
}

func (filter ImageFilter) Has(key string) bool {
	key = strings.ToLower(key)
	_, ok := filter[key]
	return ok
}

type ImageSummary struct {
	ID         string    `json:"id"`
	Repository string    `json:"repository"`
	Tag        string    `json:"tag"`
	CreatedAt  time.Time `json:"created_at"`
	Size       int64     `json:"size"`
}

func (sum ImageSummary) Name() string {
	return fmt.Sprintf("%s:%s", sum.Repository, sum.Tag)
}

func (sum ImageSummary) Language() string {
	lang := ""
	parts := strings.Split(sum.Repository, "/")
	if len(parts) > 1 {
		lang = parts[1]
	}
	return lang
}

func (sum ImageSummary) OperatingSystem() string {
	parts := strings.Split(sum.Tag, "-")
	return parts[0]
}

func (sum ImageSummary) Architecture() string {
	arch := ""
	parts := strings.Split(sum.Tag, "-")
	if len(parts) > 1 {
		arch = parts[1]
	}
	return arch
}
