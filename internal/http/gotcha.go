package http

import (
	"moderari/internal/config"
	"time"

	"github.com/fjah/gotcha"
)

// Gotcha is the main gotcha handler.
var Gotcha *gotcha.Server

// Serve starts the gotcha server. This listens on `address` in the config.
func Serve() error {
	// TODO: Implement a renderer.
	Gotcha = &gotcha.Server{
		Address: config.C.Address,
		// TODO: Make this customisable.
		Timeout: 1 * time.Hour,
	}
	return Gotcha.Serve()
}
