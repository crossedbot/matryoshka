package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/crossedbot/common/golang/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/crossedbot/matryoshka/pkg/deployer"
	"github.com/crossedbot/matryoshka/pkg/runner"
)

const (
	ReadTimeout = 60 * time.Second
	StopTimeout = 30 * time.Second
)

// Controller represents an interface to a container manager.
type Controller interface {
	// CreateDeployment runs a given code payload and returns the results.
	CreateDeployment(runner.Payload) (runner.Result, error)
}

// controller implments the Controller interface.
type controller struct {
	ctx      context.Context
	cli      *client.Client
	deployer deployer.Deployer
}

// Config represents the definable configuration for the container manager.
type Config struct {
	DockerHost       string `toml:"docker_host"`
	DockerApiVersion string `toml:"DockerApiVersion"`
	DockerTimeout    int    `toml:"docker_timeout"` // in seconds
}

// control is a singleton of a Controller and can be accessed via V* function.
var control Controller
var controllerOnce sync.Once
var V1 = func() Controller {
	// initialize the controller only once
	controllerOnce.Do(func() {
		// TODO use configuration
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			panic(err)
		}
		control = New(ctx, cli)
	})
	return control
}

// New returns a new Controller.
func New(ctx context.Context, cli *client.Client) Controller {
	d := deployer.New(ctx, cli)
	return &controller{ctx, cli, d}
}

func (c *controller) CreateDeployment(payload runner.Payload) (runner.Result, error) {
	// deploy container for the payload's programming language
	id, err := c.deploy(
		payload.Language,
		payload.OperatingSystem,
		payload.Architecture,
	)
	if err != nil {
		return runner.Result{}, err
	}
	defer c.deployer.Stop(id, StopTimeout)
	// Wait until the container stops, read all output, and parse it as a
	// result.
	output := make(chan runner.Result)
	go func() {
		b, err := c.deployer.WaitAndRead(
			id,
			container.WaitConditionNotRunning,
		)
		if err != nil {
			output <- runner.Result{Error: err.Error()}
		}
		lines := bytes.Split(b, []byte("\n"))
		for _, line := range lines {
			var result runner.Result
			err := json.Unmarshal(line, &result)
			if err == nil && (runner.Result{}) != result {
				output <- result
			}
		}
	}()
	// Write the code payload to the container
	if err := c.write(id, payload); err != nil {
		return runner.Result{}, err
	}
	// Read and return results
	return read(output)
}

// deploy deploys the appropriate container for the given language and returns
// the container's identifier.
func (c *controller) deploy(lang, os, arch string) (string, error) {
	image, err := deployer.GetImage(lang, os, arch)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("Deploying  image \"%s\"", image))
	return c.deployer.Deploy(image)
}

// write writes the given payload to the given container.
func (c *controller) write(id string, payload runner.Payload) error {
	input, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.deployer.Write(id, input)
}

// read reads all results from the output.
func read(output <-chan runner.Result) (runner.Result, error) {
	for {
		select {
		case <-time.After(ReadTimeout):
			// return if we have exceeded the timeout
			return runner.Result{}, fmt.Errorf(
				"read timeout exceeded (%s)",
				ReadTimeout.String(),
			)
		case result := <-output:
			// just grab the first result
			return result, nil
		}
	}
}
