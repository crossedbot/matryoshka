package controller

import (
	"net/http"

	"github.com/crossedbot/common/golang/server"
)

// Route represents a Deployer HTTP API route.
type Route struct {
	Handler          server.Handler
	Method           string
	Path             string
	ResponseSettings []server.ResponseSetting
}

// Routes is list of routes of the Deployer HTTP API.
var Routes = []Route{
	// Deploy and run container
	Route{
		CreateContainer,
		http.MethodPost,
		"/deployer/run",
		[]server.ResponseSetting{},
	},
}
