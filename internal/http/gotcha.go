package http

import (
	"github.com/fjah/gotcha"
	"moderari/internal/config"
	"time"
)

// Gotcha is the main gotcha handler.
var Gotcha *gotcha.Server

func Serve() error {
	// TODO: Implement a renderer.
	Gotcha = &gotcha.Server{
		Address: config.C.Address,
		// TODO: Make this customisable.
		Timeout: 30 * time.Hour,
	}
	return Gotcha.Serve()
}
