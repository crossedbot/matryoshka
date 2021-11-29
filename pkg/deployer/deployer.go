package deployer

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/crossedbot/matryoshka/pkg/deployer/models"
)

var (
	ROOT_IMAGE_NAME = "matryoshka"
)

// Deployer represents a container manager that is capable of deploying a given
// image name.
type Deployer interface {
	// Deploy deploys a given image as a container, returning the container
	// identifier on success.
	Deploy(image string) (string, error)

	// Attach attaches to the given container and returns the hijacked
	// connection.
	Attach(containerId string) (types.HijackedResponse, error)

	// Write writes the given data to the given container's STDIN.
	Write(containerId string, data []byte) error

	// WaitAndRead attaches to the given container, waits on the given wait
	// condition, and returns all output read.
	WaitAndRead(containerId string, condition container.WaitCondition) ([]byte, error)

	// Stop stops the given container, forcefully stopping the container
	// when the given timeout expires.
	Stop(containerId string, timeout time.Duration) error

	// ListImages returns a list of images that meet the given filter
	ListImages(filter models.ImageFilter) ([]models.ImageSummary, error)
}

// deployer implements the Deployer interface.
type deployer struct {
	*client.Client
	ctx context.Context
}

// New returns a new Deployer.
func New(ctx context.Context, cli *client.Client) Deployer {
	return &deployer{cli, ctx}
}

func (d *deployer) Deploy(image string) (string, error) {
	// Create the container, start it, and return its identifer
	resp, err := d.ContainerCreate(d.ctx, &container.Config{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Tty:          true,
		OpenStdin:    true,
		Image:        image,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}
	err = d.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d *deployer) Attach(containerId string) (types.HijackedResponse, error) {
	// Hijack the streams of the given container
	return d.ContainerAttach(
		d.ctx, containerId,
		types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		},
	)
}

func (d *deployer) Write(containerId string, data []byte) error {
	// Attach to the container and write the data to STDIN
	hijack, err := d.Attach(containerId)
	if err != nil {
		return err
	}
	defer hijack.Close()
	// Make sure to add a newline to delimit the end of the data stream
	_, err = fmt.Fprintf(hijack.Conn, string(data)+"\n")
	return err
}

func (d *deployer) WaitAndRead(containerId string, condition container.WaitCondition) ([]byte, error) {
	// Attach to the container
	hijack, err := d.Attach(containerId)
	if err != nil {
		return nil, err
	}
	defer hijack.Close()
	// Wait until condition is met, check status, and return all output
	dockerStatus, errStatus := d.ContainerWait(d.ctx, containerId, condition)
	select {
	case err := <-errStatus:
		if err != nil {
			return nil, err
		}
	case <-dockerStatus:
	}
	return ioutil.ReadAll(hijack.Reader)
}

func (d *deployer) Stop(containerId string, timeout time.Duration) error {
	return d.ContainerStop(d.ctx, containerId, &timeout)
}

func (d *deployer) ListImages(filter models.ImageFilter) ([]models.ImageSummary, error) {
	lang := "*"
	if v := filter.Get("language"); len(v) > 0 {
		lang = v
	}
	opSys := "*"
	if v := filter.Get("operating_system"); len(v) > 0 {
		opSys = v
	}
	arch := "*"
	if v := filter.Get("architecture"); len(v) > 0 {
		arch = v
	}
	reference := fmt.Sprintf(
		"%s/%s:%s-%s",
		ROOT_IMAGE_NAME, lang, opSys, arch,
	)
	filters := filters.NewArgs(filters.KeyValuePair{
		Key:   "reference",
		Value: reference,
	})
	imageSummaries, err := d.ImageList(
		d.ctx,
		types.ImageListOptions{Filters: filters},
	)
	if err != nil {
		return []models.ImageSummary{}, err
	}
	var images []models.ImageSummary
	for _, summary := range imageSummaries {
		id := summary.ID
		if strings.Contains(summary.ID, ":") {
			id = strings.Split(summary.ID, ":")[1]
		}
		repo := ""
		tag := ""
		if len(summary.RepoTags) > 0 {
			repoTagParts := strings.Split(summary.RepoTags[0], ":")
			repo = repoTagParts[0]
			if len(repoTagParts) > 1 {
				tag = repoTagParts[1]
			}
		}
		images = append(images, models.ImageSummary{
			ID:         id,
			Repository: repo,
			Tag:        tag,
			CreatedAt:  time.Unix(summary.Created, 0),
			Size:       summary.Size,
		})
	}
	return images, nil
}
