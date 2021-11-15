package deployer

import (
	"fmt"
	"strings"
)

// LanguageImages is a map of Docker image names and its set label.
var LanguageImages = map[string]string{
	"c":  "matryoshka/c",
	"go": "matryoshka/golang",
}

// GetImage returns the image name for a given language label.
func GetImage(lang string) (string, error) {
	lang = strings.ToLower(lang)
	if image, ok := LanguageImages[lang]; ok {
		return image, nil
	}
	return "", fmt.Errorf("unkown language \"%s\"", lang)
}
