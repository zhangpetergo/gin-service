// Package web contains a small web framework extension.
package web

import (
	"github.com/gin-gonic/gin"
	"os"
	"syscall"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*gin.Engine
	shutdown chan os.Signal
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, middleware ...gin.HandlerFunc) *App {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware...)

	return &App{
		Engine:   app,
		shutdown: shutdown,
	}
}
