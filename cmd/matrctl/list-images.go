package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/subcommands"

	"github.com/crossedbot/matryoshka/pkg/deployer/controller"
	"github.com/crossedbot/matryoshka/pkg/deployer/models"
)

type listImagesCmd struct {
	language     string
	os           string
	arch         string
	outputFormat OutputFormat
}

func (lc *listImagesCmd) Name() string {
	return "list-images"
}

func (lc *listImagesCmd) Synopsis() string {
	return `List available images for given programming language, operating
			 system, or architecture`
}

func (lc *listImagesCmd) Usage() string {
	return `list-images [-language <language>] [-os <os>] [-arch <arch>]
            [-output-format <fmt>]
  List available images.
`
}

func (lc *listImagesCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&lc.language, "language", "", "Programming language to filter images by")
	f.StringVar(&lc.os, "os", "", "Operating system to filter images by")
	f.StringVar(&lc.arch, "arch", "", "Architecture to filter images by")
	f.Var(&lc.outputFormat, "output-format", "Set the output format, plain|json (default \"plain\")")
}

func (lc *listImagesCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	images, err := controller.V1().ListImages(lc.language, lc.os, lc.arch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", lc.Name(), err)
		return subcommands.ExitFailure
	}
	outputFormat := lc.outputFormat
	if outputFormat.String() == "" {
		outputFormat = OutputFormatPlain
	}
	printImages(outputFormat, images)
	return subcommands.ExitSuccess
}

func printImages(format OutputFormat, images []models.ImageSummary) {
	switch format {
	case OutputFormatPlain:
		printImagesPlain(images)
	case OutputFormatJson:
		printImagesJson(images)
	}
}

func formatSize(sz int64) string {
	const unit = 1000
	if sz < unit {
		return fmt.Sprintf("%d B", sz)
	}
	div, exp := int64(unit), 0
	for n := sz / unit; n >= unit; n = n / unit {
		div = div * unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(sz)/float64(div), "kMGTPE"[exp])
}

func printImagesPlain(images []models.ImageSummary) {
	for _, image := range images {
		fmt.Printf("ID=%s\n", image.ID[:12])
		fmt.Printf("Repository=%s\n", image.Repository)
		fmt.Printf("Tag=%s\n", image.Tag)
		fmt.Printf("CreatedAt=%s\n",
			image.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Size=%s\n", formatSize(image.Size))
		fmt.Println()
	}
}

func printImagesJson(images []models.ImageSummary) {
	b, err := json.MarshalIndent(images, "", "    ")
	if err == nil {
		fmt.Fprintf(os.Stdout, string(b))
	}
}
