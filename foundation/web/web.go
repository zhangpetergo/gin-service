// Package web contains a small web framework extension.
package web

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, v ...any)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*gin.Engine
	log Logger
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log Logger, middleware ...gin.HandlerFunc) *App {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware...)

	return &App{
		Engine: app,
		log:    log,
	}
}
