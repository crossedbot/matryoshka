package deployer

import (
	"fmt"
	"strings"
)

// TODO leverage the docker API to pull in image information

// LanguageImages is a map of Docker image names and its set label.
var LanguageImages = map[string]string{
	"c":  "matryoshka/c",
	"go": "matryoshka/golang",
}

var ImageTags = map[string][]string{
	"matryoshka/c": []string{
		"debian-amd64",
		"debian-arm64",
	},
	"matryoshka/golang": []string{
		"debian-amd64",
		"debian-arm64",
	},
}

// GetImage returns the image name for a given language label.
func GetImage(lang, os, arch string) (string, error) {
	lang = strings.ToLower(lang)
	os = strings.ToLower(os)
	arch = strings.ToLower(arch)
	tag := fmt.Sprintf("%s-%s", os, arch)
	image, ok := LanguageImages[lang]
	if !ok {
		return "", fmt.Errorf("unkown language \"%s\"", lang)
	}
	imageTags, ok := ImageTags[image]
	if !ok {
		return "", fmt.Errorf(
			"failed to find tags for image \"%s\"",
			image,
		)
	}
	for _, imageTag := range imageTags {
		if imageTag == tag {
			return fmt.Sprintf("%s:%s", image, tag), nil
		}
	}
	return "", fmt.Errorf(
		"failed to find image for OS \"%s\" and architecture \"%s\"",
		os, arch,
	)
}
